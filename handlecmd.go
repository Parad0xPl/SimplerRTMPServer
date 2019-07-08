package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
	"net"
)

func handleCmdConnect(packet Packet, c net.Conn, cmd utils.Command) {
	head, body := Create.PCMWindowAckSize(packet.ctx.ServerWindowAcknowledgement)
	sendPacket(c, packet.ctx, head, body)
	head, body = Create.PCMSetPeerBandwitdh(packet.ctx.PeerBandwidth, 1)
	sendPacket(c, packet.ctx, head, body)
	streamID := streamsMan.createStream()
	head, body = Create.UCMStreamBegin(int(streamID))
	sendPacket(c, packet.ctx, head, body)
	head, body = Create.resultMessage(int(cmd.TransactionID), nil, amf0.Undefined{})
	sendPacket(c, packet.ctx, head, body)
}

// CMDPlayParameters _
type CMDPlayParameters struct {
	StreamName string
	Start      float64
	Duration   float64
	Reset      bool
}

func parseCMDPlayParameters(raw []interface{}) (CMDPlayParameters, error) {
	StreamName, ok := raw[0].(string)
	if !ok {
		return CMDPlayParameters{}, errors.New("StreamName is not a string")
	}

	var Start, Duration float64
	var Reset bool

	length := len(raw)
	if length >= 2 {
		Start, ok = raw[1].(float64)
		if !ok {
			return CMDPlayParameters{}, errors.New("Start is not a number")
		}
	}

	if length >= 3 {
		Duration, ok = raw[2].(float64)
		if !ok {
			return CMDPlayParameters{}, errors.New("Duration is not a number")
		}
	}

	if length >= 4 {
		Reset, ok = raw[3].(bool)
		if !ok {
			return CMDPlayParameters{}, errors.New("Reset is not a bool")
		}
	}

	return CMDPlayParameters{
		StreamName,
		Start,
		Duration,
		Reset,
	}, nil
}

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
		handleCmdConnect(packet, c, command)
	case "releaseStream":
		name, ok := raw[3].(string)
		if !ok {
			return nil
		}
		fmt.Println("Release Stream:", name)
	case "FCPublish": // Dont know where this came from
		// Founded in obs communication
		name, ok := raw[3].(string)
		if !ok {
			return nil
		}
		fmt.Println("FCPublish:", name)
	case "createStream":
		streamID := streamsMan.createStream()
		head, body := Create.resultMessage(int(command.TransactionID), nil, streamID)
		sendPacket(c, packet.ctx, head, body)

	case "play":
		_, err := parseCMDPlayParameters(raw[3:])
		if err != nil {
			return err
		}
		head, body := Create.PCMSetChunkSize(packet.ctx.ChunkSize)
		sendPacket(c, packet.ctx, head, body)
	case "play2":

	}

	return nil
}
