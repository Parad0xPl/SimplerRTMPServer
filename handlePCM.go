package main

import (
	"errors"
	"fmt"
	"net"
)

// handlePCM Handle 'Protocol Control Messages'
func handlePCM(packet Packet, c net.Conn) error {
	head := packet.header
	ctx := packet.ctx
	if head.TypeID == 1 {
		fmt.Println("Get 'Set Chunk Size'")
		var patt byte = 1
		patt <<= 7
		if patt&packet.data[0] != 0 {
			return errors.New("Wrong chunk size")
		}
		ctx.ChunkSize = readInt(packet.data)
	} else if head.TypeID == 2 {
		fmt.Println("Get 'Abort Message'")
		// TODO Abort message
	} else if head.TypeID == 3 {
		fmt.Println("Get 'Acknowledgement'")
		// TODO
	} else if head.TypeID == 5 {
		fmt.Println("Get 'Window Acknowledgement Size'")
		// TODO
	} else if head.TypeID == 6 {
		fmt.Println("Get 'Set Peer Bandwidth'")
		// TODO
	}
	return nil
}
