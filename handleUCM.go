package main

import (
	"net"
)

// handleUCM Handle 'User Control Messages'
func handleUCM(packet Packet, c net.Conn) error {
	// WHOLE TODO [PAGE 23]
	// eventType := utils.ReadInt(packet.data.bytes[0:1])
	return nil
}
