package handlers

import (
	"SimpleRTMPServer/amf0"
	"log"
)

func AMF0cmd(recvPacket ReceivedPacket) error {
	log.Println("AMF0 command")
	err, parsed := amf0.Read(recvPacket.Data, len(recvPacket.Data))
	if err != nil {
		return err
	}
	return handleCmd(recvPacket, parsed)
}
