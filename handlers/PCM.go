package handlers

import (
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
)

// handlePCM Handle 'Protocol Control Messages'
func PCM(recvPacket ReceivedPacket) error {
	head := recvPacket.Header
	ctx := recvPacket.CTX
	switch head.MessageTypeID {
	case 1:
		fmt.Println("Get 'Set Chunk Size'")
		var patt byte = 1
		patt <<= 7
		if patt&recvPacket.Data[0] != 0 {
			return errors.New("Wrong chunk size")
		}
		ctx.ChunkSize = utils.ReadInt(recvPacket.Data)
	case 2:
		fmt.Println("Get 'Abort Message'")
		// TODO Abort message
	case 3:
		fmt.Println("Get 'Acknowledgement'")
		seqNum := utils.ReadInt(recvPacket.Data[0:4])
		fmt.Println("Sequence number:", seqNum)
	case 4:
		err := handleUCM(recvPacket)
		return err
	case 5:
		fmt.Println("Get 'Window Acknowledgement Size'")
		recvPacket.CTX.ClientWindowAcknowledgement = utils.ReadInt(recvPacket.Data[0:4])
		fmt.Println("Client WinAck:", recvPacket.CTX.ClientWindowAcknowledgement)
	case 6:
		fmt.Println("Get 'Set Peer Bandwidth'")
		typ := utils.ReadInt(recvPacket.Data[4:5])
		if !(typ == 2 && recvPacket.CTX.PeerBandwidthType == 1) {
			recvPacket.CTX.PeerBandwidth = utils.ReadInt(recvPacket.Data[0:4])
			recvPacket.CTX.PeerBandwidthType = typ % 2
		}
	}
	return nil
}
