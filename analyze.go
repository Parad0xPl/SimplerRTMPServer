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

func getType(p ReceivedPacket) (string, bool) {
	if p.Header.ChunkStreamID == 2 && p.Header.MessageStreamID == 0 {
		if p.Header.MessageTypeID == 4 {
			return "UCM", false
		} else {
			return "PCM", false
		}
	} else {
		switch p.Header.MessageTypeID {
		case 8:
			return "Audio Data", false
		case 9:
			return "Video Data", false
		case 15:
			return "AMF3 Data", false
		case 17:
			return "AMF3 Command", false
		case 18:
			return "AMF0 Data", true
		case 20:
			return "AMF0 Command", true
		}
	}
	return "undefined", true
}

func analyze(fn string) {
	fmt.Printf("Analyzing %s file.\n", fn)
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

	ctx := initCTX(&conn)
	defer ctx.Clear()

	index := 1
	for {
		fmt.Printf("[%d] Position %d\n", index, conn.AmountRead())
		header, bytes, err := ctx.ReadPacket()

		p := ReceivedPacket{
			&ctx,
			&header,
			bytes,
		}

		typeOf, ifPrintData := getType(p)

		fmt.Printf("[%d] File read %.2f%%\n", index, conn.Percent()*100)
		fmt.Printf("[%d] Header: %s\n", index, header)
		fmt.Printf("[%d] Type: %s\n", index, typeOf)
		fmt.Printf("[%d] DataLen: %d\n", index, len(bytes))
		if ifPrintData {
			fmt.Printf("[%d] Data: \n%s\n", index, formatByteSlice(bytes))
		}

		if err != nil {
			fmt.Println("Can't read packet\n", err)
			return
		}

		if header.ChunkStreamID == 2 && header.MessageStreamID == 0 {
			err := handlePCM(p)
			if err != nil {
				fmt.Println("Problem with PCM", err)
			}
		}

		fmt.Println()
		index++
	}
}
