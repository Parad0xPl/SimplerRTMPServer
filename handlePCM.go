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
		tmp := make([]byte, 4)
		var patt byte = 1
		patt <<= 7
		if patt&tmp[0] != 0 {
			return errors.New("Wrong chunk size")
		}
		_, err := c.Read(tmp)
		if err != nil {
			return err
		}
		ctx.ChunkSize = readInt(tmp)
	} else if head.TypeID == 2 {
		fmt.Println("Get 'Abort Message'")
		tmp := make([]byte, 4)
		_, err := c.Read(tmp)
		if err != nil {
			return err
		}
		// TODO Abort message
	} else if head.TypeID == 3 {
		fmt.Println("Get 'Acknowledgement'")
		tmp := make([]byte, 4)
		_, err := c.Read(tmp)
		if err != nil {
			return err
		}
		// TODO
	} else if head.TypeID == 5 {
		fmt.Println("Get 'Window Acknowledgement Size'")
		tmp := make([]byte, 4)
		_, err := c.Read(tmp)
		if err != nil {
			return err
		}
		// TODO
	} else if head.TypeID == 6 {
		fmt.Println("Get 'Set Peer Bandwidth'")
		tmp := make([]byte, 5)
		_, err := c.Read(tmp)
		if err != nil {
			return err
		}
		// TODO
	}
	return nil
}
