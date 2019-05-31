package main

import (
	"SimpleRTMPServer/amf0"
	"errors"
	"log"
	"net"
)

func handleAMF0cmd(packet Packet, c net.Conn) error {
	log.Println("AMF0 command")
	parsed := amf0.Read(packet.data)
	command, ok := parsed[0].(string)
	if !ok {
		return errors.New("Wrong format")
	}
	switch command {
	case "connect":
		if !(len(parsed) >= 3) {
			return errors.New("Insufficient number of parameters")
		}
		transactionID, ok := parsed[1].(float64)
		if !ok {
			return errors.New("Wrong format")
		}
		if transactionID != 1 {
			log.Println("Transcation ID should equal 1")
		}
		commandObjectRaw, ok := parsed[2].(map[string]interface{})
		if !ok {
			return errors.New("Wrong format")
		}
		parseCommandObject(commandObjectRaw)
	}

	return nil
}
