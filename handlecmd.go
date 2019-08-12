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
	packet.CTX.SendPacket(pkt)
}

// Handlers

func handleCmdConnect(packet ReceivedPacket, cmd utils.Command) error {
	app, ok := cmd.CMDObject["app"].(string)
	if !ok {
		return errors.New("Channel name need to specified")
	}
	if !serverInstance.checkChannel(app) {
		return errors.New("Channel doesn't exists")
	}
	packet.CTX.Channel = serverInstance.Channels[app]

	pkt := Create.PCMWindowAckSize(packet.CTX.ServerWindowAcknowledgement)
	packet.CTX.SendPacket(pkt)
	pkt = Create.PAMSetPeerBandwidth(packet.CTX.PeerBandwidth, 1)
	packet.CTX.SendPacket(pkt)
	streamID := serverInstance.createStream()
	pkt = Create.UCMStreamBegin(int(streamID))
	packet.CTX.SendPacket(pkt)
	pkt = Create.resultMessage(int(cmd.TransactionID), nil, amf0.Undefined{})
	packet.CTX.SendPacket(pkt)
	return nil
}

func handleCmdCreateStream(packet ReceivedPacket, cmd utils.Command) {
	streamID := serverInstance.createStream()
	packet.CTX.NetStreamID = streamID
	log.Println("Create NetStream:", streamID)
	// streamID casted to uint as int create some strange behaviour
	// TODO FIX: streamID is not added to result massage as uint isn't supported
	// in amf0 parser. Didn't notice it, when fixed i couldn't connect with OBS
	pkt := Create.resultMessage(int(cmd.TransactionID), nil, uint(streamID))
	packet.CTX.SendPacket(pkt)
}

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

func handleCmdPlay(packet ReceivedPacket, cmd utils.Command, params CMDPlayParameters) error {
	pkt := Create.PCMSetChunkSize(packet.CTX.ChunkSize)
	packet.CTX.SendPacket(pkt)
	pkt = Create.UCMStreamIsRecorded(int(packet.CTX.ChunkStreamID))
	packet.CTX.SendPacket(pkt)
	pkt = Create.UCMStreamBegin(int(packet.CTX.ChunkStreamID))
	packet.CTX.SendPacket(pkt)
	pkt = Create.onStatusMessage("status", "NetStream.Play.Start", "Play stream")
	if params.Reset {
		packet.CTX.SendPacket(pkt)
		pkt = Create.onStatusMessage("status", "NetStream.Play.Reset", "Reset stream")
	}
	if packet.CTX.Channel.Metadata != nil {
		pkt = Create.amf0Data([]interface{}{
			"onMetaData",
			amf0.CreateECMAArray(packet.CTX.Channel.Metadata),
		})
		packet.CTX.SendPacket(pkt)
		fmt.Println("Metadata has been send")
	}
	packet.CTX.Channel.Subscribe(packet.CTX)
	return nil
}

type cmdPublishParameters struct {
	Name    string
	PubType string
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
		Name:    name,
		PubType: pubtype,
	}, nil
}

func handleCmdPublish(packet ReceivedPacket, cmd utils.Command, params cmdPublishParameters) error {
	if packet.CTX.Channel.Verify(params.Name) {
		return errors.New("Key is incorrect")
	}

	packet.CTX.AudioHandler = createAudioHandler(packet.CTX)
	packet.CTX.VideoHandler = createVideoHandler(packet.CTX)

	pkt := Create.onStatusMessage("status", "NetStream.Publish.Start", "Started publishing stream")
	packet.CTX.SendPacket(pkt)
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
			packet.CTX.Properties = &cmdObj
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
