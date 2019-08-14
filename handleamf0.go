package main

import (
	"SimpleRTMPServer/amf0"
	"log"
)

func handleAMF0cmd(packet ReceivedPacket) error {
	log.Println("AMF0 command")
	err, parsed := amf0.Read(packet.Data, len(packet.Data))
	if err != nil {
		return err
	}
	return handleCmd(packet, parsed)
}
