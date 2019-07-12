package main

import (
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
)

// handlePCM Handle 'Protocol Control Messages'
func handlePCM(packet ReceivedPacket) error {
	head := packet.header
	ctx := packet.ctx
	switch head.TypeID {
	case 1:
		fmt.Println("Get 'Set Chunk Size'")
		var patt byte = 1
		patt <<= 7
		if patt&packet.data[0] != 0 {
			return errors.New("Wrong chunk size")
		}
		ctx.ChunkSize = utils.ReadInt(packet.data)
	case 2:
		fmt.Println("Get 'Abort Message'")
		// TODO Abort message
	case 3:
		fmt.Println("Get 'Acknowledgement'")
		seqnum := utils.ReadInt(packet.data[0:4])
		fmt.Println("Sequence number:", seqnum)
	case 5:
		fmt.Println("Get 'Window Acknowledgement Size'")
		packet.ctx.ClientWindowAcknowledgement = utils.ReadInt(packet.data[0:4])
		fmt.Println("Client WinAck:", packet.ctx.ClientWindowAcknowledgement)
	case 6:
		fmt.Println("Get 'Set Peer Bandwidth'")
		typ := utils.ReadInt(packet.data[4:5])
		if !(typ == 2 && packet.ctx.PeerBandwidthType == 1) {
			packet.ctx.PeerBandwidth = utils.ReadInt(packet.data[0:4])
			packet.ctx.PeerBandwidthType = typ % 2
		}
	}
	return nil
}
