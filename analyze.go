package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/handlers"
	"SimpleRTMPServer/utils"
	"fmt"
	"strings"
)

func formatByteSlice(slc []byte) string {
	target := make([]byte, len(slc)*3)
	index := 0
	for _, v := range slc {
		target[index] = utils.ToHex(uint8(v) >> 4)
		index++
		target[index] = utils.ToHex(uint8(v) & ((1 << 4) - 1))
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

func getType(p handlers.ReceivedPacket) (string, bool) {
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

func Indent(builder *strings.Builder, indent int) {
	const indentString = "  "
	for i := 0; i < indent; i++ {
		builder.WriteString(indentString)
	}
}

func CreateString(builder *strings.Builder, data interface{}, indent int) {
	switch typed := data.(type) {
	case []interface{}:
		Indent(builder, indent)
		builder.WriteString("[\n")

		for _, v := range typed {
			CreateString(builder, v, indent+1)
		}

		Indent(builder, indent)
		builder.WriteString("]\n")
	case map[string]interface{}:
		Indent(builder, indent)
		builder.WriteString("{\n")
		for key, v := range typed {
			Indent(builder, indent+1)
			builder.WriteString(fmt.Sprintf("[%s]: %v\n", key, v))
		}
		Indent(builder, indent)
		builder.WriteString("}\n")
	default:
		Indent(builder, indent)
		builder.WriteString(fmt.Sprintln(typed))
	}
}

func Print(data interface{}) {
	builder := new(strings.Builder)
	CreateString(builder, data, 0)
	fmt.Print(builder.String())
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

	ctx := NewCTX(&conn)
	defer ctx.Clear()

	index := 1
	for {
		amountRead := conn.AmountRead()
		header, bytes, err := ctx.ReadPacket()

		fmt.Printf("[%d] Position %d\n", index, amountRead)

		p := handlers.ReceivedPacket{
			ctx,
			serverInstance,
			&header,
			bytes,
		}

		typeOf, ifPrintData := getType(p)

		fmt.Printf("[%d] File read %.2f%%\n", index, conn.Percent()*100)
		fmt.Printf("[%d] Header: %s\n", index, header)
		fmt.Printf("[%d] Type: %s\n", index, typeOf)
		fmt.Printf("[%d] DataLen: %d\n", index, len(bytes))

		if err != nil {
			if header.ChunkStreamID == 0 {
				fmt.Println("End of file")
				return
			}
			fmt.Println("Can't read packet\n", err)
			return
		}

		if header.ChunkStreamID == 2 && header.MessageStreamID == 0 {
			err := handlers.PCM(p)
			if err != nil {
				fmt.Println("Problem with PCM", err)
			}
		} else {
			switch header.MessageTypeID {
			case 18:
				err, decoded := amf0.Read(bytes, len(bytes))
				if err != nil {
					fmt.Println("[ERR]AMF0 Data corruption:", err)
					fmt.Println("[ERR]Part")
					Print(decoded)
				} else {
					Print(decoded)
				}
			case 20:
				err, decoded := amf0.Read(bytes, len(bytes))
				if err != nil {
					fmt.Println("[ERR]AMF0 Command corruption:", err)
					fmt.Println("[ERR]Part")
					Print(decoded)
				} else {
					Print(decoded)
				}
			default:
				if ifPrintData {
					fmt.Printf("[%d] Data: \n%s\n", index, formatByteSlice(bytes))
				}
			}
		}

		fmt.Println()
		index++
	}
}
