package main

import (
	"SimpleRTMPServer/utils"
	"net"
)

// handleUCM Handle 'User Control Messages'
func handleUCM(packet Packet, c net.Conn) error {
	// WHOLE TODO [PAGE 23]
	eventType := utils.ReadInt(packet.data[0:2])

	switch eventType {
	case 0:
		// Stream Begin
		/*StreamId=*/
		utils.ReadInt(packet.data[2:6])
	case 1:
		// Stream EOF
		/*StreamId=*/
		utils.ReadInt(packet.data[2:6])
	case 2:
		// StreamDry
		/*StreamId=*/
		utils.ReadInt(packet.data[2:6])
	case 3:
		// SetBuffer Length
		/*StreamId=*/
		utils.ReadInt(packet.data[2:6])
		/*BufferLength=*/
		utils.ReadInt(packet.data[6:10])
	case 4:
		// StreamIs Recorded
		/*StreamId=*/
		utils.ReadInt(packet.data[2:6])
	case 6:
		// PingRequest
		pkt := Create.UCMPingResponse(packet.data[2:6])
		sendPacket(c, packet.ctx, pkt)
	case 7:
		// PingResponse
		// Only client
	}

	return nil
}
