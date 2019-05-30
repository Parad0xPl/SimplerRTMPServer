package main

import (
	"SimpleRTMPServer/amf0"
	"errors"
	"fmt"
	"log"
	"net"
)

func handleAMF0cmd(packet Packet, c net.Conn) error {
	log.Println("AMF0 command")
	var command string
	var n int
	if packet.data[0] != 2 {
		return errors.New("Command doesn't start with string")
	}
	i := 0
	i++
	command, n = amf0.ReadString(packet.data[i:])
	i += n
	log.Println("Command recived:", command)

	i++
	transactionID, n := amf0.ReadNumber(packet.data[i:])
	i += n
	log.Println("TransactionId:", transactionID)
	options, n := amf0.ReadObject(packet.data[13+len(command):])
	for i, n := range options {
		fmt.Println(i, n)
	}
	return nil
}
