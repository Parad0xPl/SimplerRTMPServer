package handlers

import (
	"SimpleRTMPServer/packet"
	"SimpleRTMPServer/utils"
)

// handleUCM Handle 'User Control Messages'
func handleUCM(pkt ReceivedPacket) error {
	// WHOLE TODO [PAGE 23]
	eventType := utils.ReadInt(pkt.Data[0:2])

	switch eventType {
	case 0:
		// Stream Begin
		/*StreamId=*/
		utils.ReadInt(pkt.Data[2:6])
	case 1:
		// Stream EOF
		/*StreamId=*/
		utils.ReadInt(pkt.Data[2:6])
	case 2:
		// StreamDry
		/*StreamId=*/
		utils.ReadInt(pkt.Data[2:6])
	case 3:
		// SetBuffer Length
		/*StreamId=*/
		utils.ReadInt(pkt.Data[2:6])
		/*BufferLength=*/
		utils.ReadInt(pkt.Data[6:10])
	case 4:
		// StreamIs Recorded
		/*StreamId=*/
		utils.ReadInt(pkt.Data[2:6])
	case 6:
		// PingRequest
		prototype := packet.Create.UCMPingResponse(pkt.Data[2:6])
		pkt.CTX.SendPacket(prototype)
	case 7:
		// PingResponse
		// Only client
	}

	return nil
}
