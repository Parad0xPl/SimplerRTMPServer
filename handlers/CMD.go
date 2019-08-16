package handlers

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/packet"
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
	"log"
)

// sendError wrapper
func sendError(recvPacket ReceivedPacket, cmd utils.Command, desc string) {
	pkt := packet.Create.ErrorMessage(packet.CommandArgs{
		int(cmd.TransactionID),
		nil,
		map[string]interface{}{
			"code":        "error",
			"level":       "error",
			"description": desc,
		},
	})
	recvPacket.CTX.SendPacket(pkt)
}

// Handlers

func handleCmdConnect(recvPacket ReceivedPacket, cmd utils.Command) error {
	app, ok := cmd.CMDObject["app"].(string)
	if !ok {
		return errors.New("Channel name need to specified")
	}
	if recvPacket.ServerInstance.CheckChannel(app) {
		return errors.New("Channel doesn't exists")
	}
	recvPacket.CTX.Channel = recvPacket.ServerInstance.Channels[app]

	pkt := packet.Create.PCMWindowAckSize(recvPacket.CTX.ServerWindowAcknowledgement)
	recvPacket.CTX.SendPacket(pkt)
	pkt = packet.Create.PCMSetPeerBandwidth(recvPacket.CTX.PeerBandwidth, 1)
	recvPacket.CTX.SendPacket(pkt)
	streamID := recvPacket.ServerInstance.CreateStream()
	pkt = packet.Create.UCMStreamBegin(int(streamID))
	recvPacket.CTX.SendPacket(pkt)
	pkt = packet.Create.ResultMessage(packet.CommandArgs{
		int(cmd.TransactionID),
		map[string]interface{}{
			"mode":         1,
			"capabilities": 127,
			"fmsVer":       "FMS/3,5,3,888",
		},
		amf0.Undefined{},
	})
	pkt.Head.MessageStreamID = recvPacket.CTX.NetStreamID
	recvPacket.CTX.SendPacket(pkt)
	return nil
}

func handleCmdCreateStream(recvPacket ReceivedPacket, cmd utils.Command) {
	streamID := recvPacket.ServerInstance.CreateStream()
	recvPacket.CTX.NetStreamID = streamID
	log.Println("Create NetStream:", streamID)
	// streamID casted to uint as int create some strange behaviour
	pkt := packet.Create.ResultMessage(packet.CommandArgs{
		int(cmd.TransactionID),
		nil,
		uint(streamID),
	})
	recvPacket.CTX.SendPacket(pkt)
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

func handleCmdPlay(recvPacket ReceivedPacket, cmd utils.Command, params CMDPlayParameters) error {
	pkt := packet.Create.PCMSetChunkSize(recvPacket.CTX.ChunkSize)
	recvPacket.CTX.SendPacket(pkt)
	pkt = packet.Create.UCMStreamIsRecorded(int(recvPacket.CTX.ChunkStreamID))
	recvPacket.CTX.SendPacket(pkt)
	pkt = packet.Create.UCMStreamBegin(int(recvPacket.CTX.ChunkStreamID))
	recvPacket.CTX.SendPacket(pkt)
	if params.Reset {
		pkt = packet.Create.OnStatusMessage("status", "NetStream.Play.Reset", "Reset stream")
		recvPacket.CTX.SendPacket(pkt)
	}
	recvPacket.CTX.SendPacket(pkt)
	pkt = packet.Create.OnStatusMessage("status", "NetStream.Play.Start", "Play stream")
	if recvPacket.CTX.Channel.Metadata != nil {
		pkt = packet.Create.AMF0Data([]interface{}{
			"onMetaData",
			amf0.CreateECMAArray(recvPacket.CTX.Channel.Metadata),
		})
		recvPacket.CTX.SendPacket(pkt)
		fmt.Println("Metadata has been send")
	}
	recvPacket.CTX.Channel.Subscribe(recvPacket.CTX)
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

func handleCmdPublish(recvPacket ReceivedPacket, cmd utils.Command, params cmdPublishParameters) error {
	if recvPacket.CTX.Channel.Verify(params.Name) {
		return errors.New("Key is incorrect")
	}

	recvPacket.CTX.AudioHandler = createAudioHandler(recvPacket.CTX)
	recvPacket.CTX.VideoHandler = createVideoHandler(recvPacket.CTX)

	pkt := packet.Create.OnStatusMessage("status", "NetStream.Publish.Start", "Started publishing stream")
	recvPacket.CTX.SendPacket(pkt)
	return nil
}

func handleCmd(recvPacket ReceivedPacket, raw []interface{}) error {
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
			recvPacket.CTX.Properties = &cmdObj
		}
		err = handleCmdConnect(recvPacket, command)

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
		handleCmdCreateStream(recvPacket, command)

	case "play":
		params, err := parseCMDPlayParameters(raw[3:])
		if err != nil {
			break
		}
		handleCmdPlay(recvPacket, command, params)

	case "publish":
		params, err := parseCMDPublishParameters(raw[3:])
		if err != nil {
			break
		}
		err = handleCmdPublish(recvPacket, command, params)

	}

	if err != nil {
		sendError(recvPacket, command, err.Error())
	}

	return err
}
