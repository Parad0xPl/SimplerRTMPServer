package main

import (
	"SimpleRTMPServer/utils"
	"net"
)

var streamID = 3

func handleCmd(packet Packet, c net.Conn, raw []interface{}) error {
	command, err := utils.ParseCommand(raw)
	if err != nil {
		return err
	}
	switch command.Name {
	case "connect":
		cmdObj, ok := raw[2].(map[string]interface{})
		if ok {
			packet.ctx.Properties = &cmdObj
		}
		head, body := Create.PCMWindowAckSize(128)
		sendPacket(c, packet.ctx, head, body)
		head, body = Create.PCMSetPeerBandwitdh(128, 1)
		sendPacket(c, packet.ctx, head, body)
		head, body = Create.UCMStreamBegin(streamID)
		sendPacket(c, packet.ctx, head, body)
		streamID++
		head, body = Create.resultMessage(*packet.ctx.Properties, make(map[string]interface{}))
		sendPacket(c, packet.ctx, head, body)
	}

	return nil
}
