package main

import (
	"SimpleRTMPServer/utils"
	"fmt"
)

func toHex(b uint8) byte {
	if b >= 10 {
		return 'a' + b - 10
	}
	return '0' + b
}

func formatByteSlice(slc []byte) string {
	target := make([]byte, len(slc)*3)
	index := 0
	for _, v := range slc {
		target[index] = toHex(uint8(v) >> 4)
		index++
		target[index] = toHex(uint8(v) & ((1 << 4) - 1))
		index++
		if ((index/3)+1)%16 == 0 {
			target[index] = '\n'
		} else {
			target[index] = ' '
		}
		index++
	}
	return string(target)
}

func analyze(fn string) {
	fmt.Printf("Analysing %s file.\n", fn)
	conn, err := utils.OpenFileConn(fn, "")
	if err != nil {
		fmt.Println("Problem with opening file:\n", err)
		return
	}

	ver := make([]byte, 1)
	_, err = conn.Read(ver)
	if err != nil {
		fmt.Println("Problem with reading first byte\n", err)
		return
	}
	if ver[0] != 3 {
		fmt.Println("Version byte don't match\n", err)
		return
	}

	_, err = conn.Seek(2*1536, 1)
	if err != nil {
		fmt.Println("Problem occurred while seeking handshake\n", err)
		return
	}

	ctx := initCTX(conn)
	defer ctx.Clear()

	index := 1
	for {
		header, bytes, err := ctx.ReadPacket()
		fmt.Printf("[%d] Header: %s\n", index, header)
		fmt.Printf("[%d] DataLen: %d\n", index, len(bytes))
		fmt.Printf("[%d] Data: \n%s\n", index, formatByteSlice(bytes))
		if err != nil {
			fmt.Println("Can't read packet\n", err)
			return
		}

		handlePacket(ReceivedPacket{
			&ctx,
			&header,
			bytes,
		})

		index++
	}
}
