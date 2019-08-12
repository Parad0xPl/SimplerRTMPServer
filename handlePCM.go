package main

import (
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
)

// handlePCM Handle 'Protocol Control Messages'
func handlePCM(packet ReceivedPacket) error {
	head := packet.Header
	ctx := packet.CTX
	switch head.MessageTypeID {
	case 1:
		fmt.Println("Get 'Set Chunk Size'")
		var patt byte = 1
		patt <<= 7
		if patt&packet.Data[0] != 0 {
			return errors.New("Wrong chunk size")
		}
		ctx.ChunkSize = utils.ReadInt(packet.Data)
	case 2:
		fmt.Println("Get 'Abort Message'")
		// TODO Abort message
	case 3:
		fmt.Println("Get 'Acknowledgement'")
		seqNum := utils.ReadInt(packet.Data[0:4])
		fmt.Println("Sequence number:", seqNum)
	case 4:
		err := handleUCM(packet)
		return err
	case 5:
		fmt.Println("Get 'Window Acknowledgement Size'")
		packet.CTX.ClientWindowAcknowledgement = utils.ReadInt(packet.Data[0:4])
		fmt.Println("Client WinAck:", packet.CTX.ClientWindowAcknowledgement)
	case 6:
		fmt.Println("Get 'Set Peer Bandwidth'")
		typ := utils.ReadInt(packet.Data[4:5])
		if !(typ == 2 && packet.CTX.PeerBandwidthType == 1) {
			packet.CTX.PeerBandwidth = utils.ReadInt(packet.Data[0:4])
			packet.CTX.PeerBandwidthType = typ % 2
		}
	}
	return nil
}
