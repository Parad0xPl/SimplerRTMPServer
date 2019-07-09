package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
)

func handleCmdConnect(packet ReceivedPacket, cmd utils.Command) {
	pkt := Create.PCMWindowAckSize(packet.ctx.ServerWindowAcknowledgement)
	sendPacket(packet.ctx, pkt)
	pkt = Create.PCMSetPeerBandwitdh(packet.ctx.PeerBandwidth, 1)
	sendPacket(packet.ctx, pkt)
	streamID := streamsMan.createStream()
	pkt = Create.UCMStreamBegin(int(streamID))
	sendPacket(packet.ctx, pkt)
	pkt = Create.resultMessage(int(cmd.TransactionID), nil, amf0.Undefined{})
	sendPacket(packet.ctx, pkt)
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

func handleCmd(packet ReceivedPacket, raw []interface{}) error {
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
		handleCmdConnect(packet, command)
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
		packet.ctx.StreamID = streamID
		pkt := Create.resultMessage(int(command.TransactionID), nil, streamID)
		sendPacket(packet.ctx, pkt)

	case "play":
		params, err := parseCMDPlayParameters(raw[3:])
		if err != nil {
			return err
		}
		pkt := Create.PCMSetChunkSize(packet.ctx.ChunkSize)
		sendPacket(packet.ctx, pkt)
		pkt = Create.UCMStreamIsRecorded(int(packet.ctx.StreamID))
		sendPacket(packet.ctx, pkt)
		pkt = Create.UCMStreamBegin(int(packet.ctx.StreamID))
		sendPacket(packet.ctx, pkt)
		pkt = Create.onStatusMessage("status", "NetStream.Play.Start", "Play stream")
		if params.Reset {
			sendPacket(packet.ctx, pkt)
			pkt = Create.onStatusMessage("status", "NetStream.Play.Reset", "Reset stream")
		}
	case "publish":
		name, ok := raw[3].(string)
		if !ok {
			return errors.New("Publishing Name is not string")
		}

		pubtype, ok := raw[4].(string)
		if !ok {
			return errors.New("Publishing Type is not string")
		}

		if pubtype != "live" {
			return errors.New("SimplerRTMPServer doesn't support other type than livestriming")
		}

		err := streamsMan.publish(packet.ctx.StreamID, name)
		if err != nil {
			return err
		}

		pkt := Create.onStatusMessage("status", "NetStream.Publish.Start", "Started publishing stream")
		sendPacket(packet.ctx, pkt)
	}

	return nil
}
