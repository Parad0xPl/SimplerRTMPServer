package main

import (
	"SimpleRTMPServer/utils"
	"net"
)

// handleUCM Handle 'User Control Messages'
func handleUCM(packet Packet, c net.Conn) error {
	// WHOLE TODO [PAGE 23]
	eventType := utils.ReadInt(packet.data.bytes[0:2])

	switch eventType {
	case 0:
		// Stream Begin
	case 1:
		// Stream EOF
	case 2:
		// StreamDry
	case 3:
		// SetBuffer Length
	case 4:
		// StreamIs Recorded
	case 6:
		// PingRequest
		head, body := Create.UCMPingResponse(packet.data.bytes[2:6])
		sendPacket(c, packet.ctx, head, body)
	case 7:
		// PingResponse
	}

	return nil
}
