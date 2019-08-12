package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/build"
	"SimpleRTMPServer/utils"
)

type create struct{}

// Create a packet
var Create create

// PacketPrototype packet skeleton
type PacketPrototype struct {
	Head Header
	Body []byte
}

func (create) PCMSetChunkSize(size int) PacketPrototype {
	head := Header{
		MessageTypeID:   1,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.SetChunkSize(size)
	return PacketPrototype{head, body}
}

func (create) PCMAbortMessage(chunkID int) PacketPrototype {
	head := Header{
		MessageTypeID:   2,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(chunkID)
	return PacketPrototype{head, body}
}

func (create) PCMAcknowledgement(seqNumber int) PacketPrototype {
	head := Header{
		MessageTypeID:   3,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(seqNumber)
	return PacketPrototype{head, body}
}

func (create) PCMWindowAckSize(winsize int) PacketPrototype {
	head := Header{
		MessageTypeID:   5,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(winsize)
	return PacketPrototype{head, body}
}

func (create) PAMSetPeerBandwidth(windowSize, limitType int) PacketPrototype {
	head := Header{
		MessageTypeID:   6,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.SetPeerBandwidth(windowSize, limitType)
	return PacketPrototype{head, body}
}

func (create) UCMStreamBegin(streamID int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(0, utils.WriteIntBE(streamID, 4))
	return PacketPrototype{head, body}
}

func (create) UCMStreamEOF(streamID int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(1, utils.WriteIntBE(streamID, 4))
	return PacketPrototype{head, body}
}

func (create) UCMStreamDry(streamID int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(2, utils.WriteIntBE(streamID, 4))
	return PacketPrototype{head, body}
}

func (create) UCMSetBufferLength(streamID, buffLen int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	eventData := utils.Concat(utils.WriteIntBE(streamID, 4), utils.WriteIntBE(buffLen, 4))
	body, _ := build.Body.UCM(3, eventData)
	return PacketPrototype{head, body}
}

func (create) UCMStreamIsRecorded(streamID int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(4, utils.WriteIntBE(streamID, 4))
	return PacketPrototype{head, body}
}

func (create) UCMPingRequest(timestamp int) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(6, utils.WriteIntBE(timestamp, 4))
	return PacketPrototype{head, body}
}

func (create) UCMPingResponse(timestamp interface{}) PacketPrototype {
	head := Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	var body []byte
	switch v := timestamp.(type) {
	case int:
		body, _ = build.Body.UCM(7, utils.WriteIntBE(v, 4))
	case []byte:
		body, _ = build.Body.UCM(7, v)
	}
	return PacketPrototype{head, body}
}

func (create) commandMessage(raw []interface{}) PacketPrototype {
	head := Header{
		MessageTypeID: 20,
		ChunkStreamID: 3,
	}
	body := amf0.Write(raw)
	return PacketPrototype{head, body}
}

type CommandArgs struct {
	TransactionID int
	Properties    interface{}
	Information   interface{}
}

func (create) resultMessage(args CommandArgs) PacketPrototype {
	raw := make([]interface{}, 4)
	raw[0] = "_result"
	raw[1] = args.TransactionID
	raw[2] = args.Properties
	raw[3] = args.Information
	packet := Create.commandMessage(raw)
	return packet
}

func (create) errorMessage(args CommandArgs) PacketPrototype {
	raw := make([]interface{}, 4)
	raw[0] = "_error"
	raw[1] = args.TransactionID
	raw[2] = args.Properties
	raw[3] = args.Information
	packet := Create.commandMessage(raw)
	return packet
}

func (create) onStatusMessage(level, code, desc string) PacketPrototype {
	raw := make([]interface{}, 4)
	raw[0] = "onStatus"
	raw[1] = 0
	raw[2] = nil
	raw[3] = map[string]interface{}{
		"level":       level,
		"code":        code,
		"description": desc,
	}
	packet := Create.commandMessage(raw)
	return packet
}

func (create) AudioData(data []byte) PacketPrototype {
	head := Header{
		MessageTypeID: 8,
	}
	return PacketPrototype{head, data}
}

func (create) VideoData(data []byte) PacketPrototype {
	head := Header{
		MessageTypeID: 9,
	}
	return PacketPrototype{head, data}
}

func (create) amf0Data(data []interface{}) PacketPrototype {
	head := Header{
		MessageTypeID: 0x12,
	}
	body := amf0.Write(data)
	return PacketPrototype{head, body}
}
