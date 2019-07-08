package main

import (
	"SimpleRTMPServer/amf0"
	"log"
	"net"
)

func handleAMF0cmd(packet Packet, c net.Conn) error {
	log.Println("AMF0 command")
	parsed := amf0.Read(packet.data)
	return handleCmd(packet, c, parsed)
}
