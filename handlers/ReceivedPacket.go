package handlers

import (
	"SimpleRTMPServer/connCTX"
	"SimpleRTMPServer/packet"
)

// ReceivedPacket wrapper for all packet related data
type ReceivedPacket struct {
	CTX            *connCTX.ConnContext
	ServerInstance *connCTX.ServerInstance
	Header         *packet.Header
	Data           []byte // Should work as long any handler won't try append to slice
}
