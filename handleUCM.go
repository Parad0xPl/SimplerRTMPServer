package main

import (
	"SimpleRTMPServer/utils"
)

// handleUCM Handle 'User Control Messages'
func handleUCM(packet ReceivedPacket) error {
	// WHOLE TODO [PAGE 23]
	eventType := utils.ReadInt(packet.Data[0:2])

	switch eventType {
	case 0:
		// Stream Begin
		/*StreamId=*/
		utils.ReadInt(packet.Data[2:6])
	case 1:
		// Stream EOF
		/*StreamId=*/
		utils.ReadInt(packet.Data[2:6])
	case 2:
		// StreamDry
		/*StreamId=*/
		utils.ReadInt(packet.Data[2:6])
	case 3:
		// SetBuffer Length
		/*StreamId=*/
		utils.ReadInt(packet.Data[2:6])
		/*BufferLength=*/
		utils.ReadInt(packet.Data[6:10])
	case 4:
		// StreamIs Recorded
		/*StreamId=*/
		utils.ReadInt(packet.Data[2:6])
	case 6:
		// PingRequest
		pkt := Create.UCMPingResponse(packet.Data[2:6])
		packet.CTX.sendPacket(pkt)
	case 7:
		// PingResponse
		// Only client
	}

	return nil
}
