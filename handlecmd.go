package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/utils"
	"fmt"
	"net"
)

var streamID = 3

func handleCmd(packet Packet, c net.Conn, raw []interface{}) error {
	command, err := utils.ParseCommand(raw)
	if err != nil {
		return err
	}
	fmt.Println("CMD name:", command.Name)
	switch command.Name {
	case "connect":
		cmdObj, ok := raw[2].(map[string]interface{})
		if ok {
			packet.ctx.Properties = &cmdObj
		}
		head, body := Create.PCMWindowAckSize(packet.ctx.ServerWindowAcknowledgement)
		sendPacket(c, packet.ctx, head, body)
		head, body = Create.PCMSetPeerBandwitdh(packet.ctx.PeerBandwidth, 1)
		sendPacket(c, packet.ctx, head, body)
		head, body = Create.UCMStreamBegin(streamID)
		sendPacket(c, packet.ctx, head, body)
		streamID++
		head, body = Create.resultMessage(int(command.TransactionID), nil, amf0.Undefined{})
		sendPacket(c, packet.ctx, head, body)
	case "releaseStream":
		name, ok := raw[3].(string)
		if !ok {
			return nil
		}
		fmt.Println("Release Stream:", name)
	case "FCPublish":
		name, ok := raw[3].(string)
		if !ok {
			return nil
		}
		fmt.Println("FCPublish:", name)
	case "createStream":
		head, body := Create.resultMessage(int(command.TransactionID), nil, streamID)
		streamID++
		sendPacket(c, packet.ctx, head, body)
	}

	return nil
}
