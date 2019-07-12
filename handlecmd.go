package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
	"log"
)

// sendError wrapper
func sendError(packet ReceivedPacket, cmd utils.Command, desc string) {
	pkt := Create.errorMessage(int(cmd.TransactionID), nil, map[string]interface{}{
		"code":        "error",
		"level":       "error",
		"description": desc,
	})
	packet.ctx.sendPacket(pkt)
}

// Handlers

func handleCmdConnect(packet ReceivedPacket, cmd utils.Command) error {
	app, ok := cmd.CMDObject["app"].(string)
	if !ok {
		return errors.New("Channel name need to specified")
	}
	if !streamsMan.checkChannel(app) {
		return errors.New("Channel doesn't exists")
	}
	packet.ctx.channel = streamsMan.channels[app]

	pkt := Create.PCMWindowAckSize(packet.ctx.ServerWindowAcknowledgement)
	packet.ctx.sendPacket(pkt)
	pkt = Create.PCMSetPeerBandwitdh(packet.ctx.PeerBandwidth, 1)
	packet.ctx.sendPacket(pkt)
	streamID := streamsMan.createStream()
	pkt = Create.UCMStreamBegin(int(streamID))
	packet.ctx.sendPacket(pkt)
	pkt = Create.resultMessage(int(cmd.TransactionID), nil, amf0.Undefined{})
	packet.ctx.sendPacket(pkt)
	return nil
}

func handleCmdCreateStream(packet ReceivedPacket, cmd utils.Command) {
	log.Println("Create to string")
	streamID := streamsMan.createStream()
	packet.ctx.StreamID = streamID
	pkt := Create.resultMessage(int(cmd.TransactionID), nil, streamID)
	packet.ctx.sendPacket(pkt)
}

type cmdPlayParameters struct {
	StreamName string
	Start      float64
	Duration   float64
	Reset      bool
}

func parseCMDPlayParameters(raw []interface{}) (cmdPlayParameters, error) {
	StreamName, ok := raw[0].(string)
	if !ok {
		return cmdPlayParameters{}, errors.New("StreamName is not a string")
	}

	var Start, Duration float64
	var Reset bool

	length := len(raw)
	if length >= 2 {
		Start, ok = raw[1].(float64)
		if !ok {
			return cmdPlayParameters{}, errors.New("Start is not a number")
		}
	}

	if length >= 3 {
		Duration, ok = raw[2].(float64)
		if !ok {
			return cmdPlayParameters{}, errors.New("Duration is not a number")
		}
	}

	if length >= 4 {
		Reset, ok = raw[3].(bool)
		if !ok {
			return cmdPlayParameters{}, errors.New("Reset is not a bool")
		}
	}

	return cmdPlayParameters{
		StreamName,
		Start,
		Duration,
		Reset,
	}, nil
}

func handleCmdPlay(packet ReceivedPacket, cmd utils.Command, params cmdPlayParameters) error {
	pkt := Create.PCMSetChunkSize(packet.ctx.ChunkSize)
	packet.ctx.sendPacket(pkt)
	pkt = Create.UCMStreamIsRecorded(int(packet.ctx.StreamID))
	packet.ctx.sendPacket(pkt)
	pkt = Create.UCMStreamBegin(int(packet.ctx.StreamID))
	packet.ctx.sendPacket(pkt)
	pkt = Create.onStatusMessage("status", "NetStream.Play.Start", "Play stream")
	if params.Reset {
		packet.ctx.sendPacket(pkt)
		pkt = Create.onStatusMessage("status", "NetStream.Play.Reset", "Reset stream")
	}
	if packet.ctx.channel.metadata != nil {
		pkt = Create.amf0Data(int(packet.ctx.StreamID), []interface{}{
			"onMetaData",
			amf0.CreateECMAArray(packet.ctx.channel.metadata),
		})
		packet.ctx.sendPacket(pkt)
		fmt.Println("Metadata has been send")
	}
	packet.ctx.channel.subscribe(packet.ctx)
	return nil
}

type cmdPublishParameters struct {
	name    string
	pubtype string
}

func parseCMDPublishParameters(raw []interface{}) (cmdPublishParameters, error) {
	name, ok := raw[0].(string)
	var err error
	if !ok {
		err = errors.New("Publishing Name is not string")
		return cmdPublishParameters{}, err
	}
	pubtype, ok := raw[1].(string)
	if !ok {
		err = errors.New("Publishing Type is not string")
		return cmdPublishParameters{}, err
	}
	if pubtype != "live" {
		err = errors.New("SimplerRTMPServer doesn't support other type than livestriming")
		return cmdPublishParameters{}, err
	}
	return cmdPublishParameters{
		name:    name,
		pubtype: pubtype,
	}, nil
}

func handleCmdPublish(packet ReceivedPacket, cmd utils.Command, params cmdPublishParameters) error {
	if packet.ctx.channel.verify(params.name) {
		return errors.New("Key is incorrect")
	}

	packet.ctx.audioHandler = createAudioHandler(packet.ctx)
	packet.ctx.videoHandler = createVideoHandler(packet.ctx)

	pkt := Create.onStatusMessage("status", "NetStream.Publish.Start", "Started publishing stream")
	packet.ctx.sendPacket(pkt)
	return nil
}

func handleCmd(packet ReceivedPacket, raw []interface{}) error {
	var err error
	err = nil

	command, err := utils.ParseCommand(raw)
	if err != nil {
		return err
	}
	fmt.Println("Command packet")
	fmt.Println("CMD name:", command.Name)

	switch command.Name {
	case "connect":
		cmdObj, ok := raw[2].(map[string]interface{})
		if ok {
			packet.ctx.Properties = &cmdObj
		}
		err = handleCmdConnect(packet, command)

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
		handleCmdCreateStream(packet, command)

	case "play":
		params, err := parseCMDPlayParameters(raw[3:])
		if err != nil {
			break
		}
		handleCmdPlay(packet, command, params)

	case "publish":
		params, err := parseCMDPublishParameters(raw[3:])
		if err != nil {
			break
		}
		err = handleCmdPublish(packet, command, params)

	}

	if err != nil {
		sendError(packet, command, err.Error())
	}

	return err
}
