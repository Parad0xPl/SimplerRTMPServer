package main

import (
	"SimpleRTMPServer/amf0"
	"log"
)

func handleAMF0cmd(packet ReceivedPacket) error {
	log.Println("AMF0 command")
	parsed := amf0.Read(packet.Data)
	return handleCmd(packet, parsed)
}
