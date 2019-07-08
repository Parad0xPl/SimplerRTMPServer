package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/build"
	"SimpleRTMPServer/utils"
)

type create struct{}

// Create a packet
var Create create

// PacketProt packet skeleton
type PacketProt struct {
	head Header
	body []byte
}

func (create) PCMSetChunkSize(size int) PacketProt {
	head := Header{
		TypeID:   1,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.PCM.SetChunkSize(size)
	return PacketProt{head, body}
}

func (create) PCMAbortMessage(chunkid int) PacketProt {
	head := Header{
		TypeID:   2,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.PCM.AbortMessage(chunkid)
	return PacketProt{head, body}
}

func (create) PCMAcknowledgement(seqnumber int) PacketProt {
	head := Header{
		TypeID:   3,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.PCM.AbortMessage(seqnumber)
	return PacketProt{head, body}
}

func (create) PCMWindowAckSize(winsize int) PacketProt {
	head := Header{
		TypeID:   5,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.PCM.AbortMessage(winsize)
	return PacketProt{head, body}
}

func (create) PCMSetPeerBandwitdh(windowsize, limittype int) PacketProt {
	head := Header{
		TypeID:   6,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.PCM.SetPeerBandwitdh(windowsize, limittype)
	return PacketProt{head, body}
}

func (create) UCMStreamBegin(streamID int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.UCM(0, utils.WriteInt(streamID, 4))
	return PacketProt{head, body}
}

func (create) UCMStreamEOF(streamID int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.UCM(1, utils.WriteInt(streamID, 4))
	return PacketProt{head, body}
}

func (create) UCMStreamDry(streamID int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.UCM(2, utils.WriteInt(streamID, 4))
	return PacketProt{head, body}
}

func (create) UCMSetBufferLength(streamID, bufflen int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	eventData := append(utils.WriteInt(streamID, 4), utils.WriteInt(bufflen, 4)...)
	body, _ := build.Body.UCM(3, eventData)
	return PacketProt{head, body}
}

func (create) UCMStreamIsRecorded(streamID int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.UCM(4, utils.WriteInt(streamID, 4))
	return PacketProt{head, body}
}

func (create) UCMPingRequest(timestamp int) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	body, _ := build.Body.UCM(6, utils.WriteInt(timestamp, 4))
	return PacketProt{head, body}
}

func (create) UCMPingResponse(timestamp interface{}) PacketProt {
	head := Header{
		TypeID:   4,
		StreamID: 0,
		ChunkID:  2,
	}
	var body []byte
	switch v := timestamp.(type) {
	case int:
		body, _ = build.Body.UCM(7, utils.WriteInt(v, 4))
	case []byte:
		body, _ = build.Body.UCM(7, v)
	}
	return PacketProt{head, body}
}

func (create) commandMessage(raw []interface{}) PacketProt {
	head := Header{
		TypeID:  20,
		ChunkID: 3,
	}
	body := amf0.Write(raw)
	return PacketProt{head, body}
}

func (create) resultMessage(transID int, props, infos interface{}) PacketProt {
	raw := make([]interface{}, 4)
	raw[0] = "_result"
	raw[1] = transID
	raw[2] = props
	raw[3] = infos
	packet := Create.commandMessage(raw)
	return packet
}

func (create) errorMessage(transID int, props, infos interface{}) PacketProt {
	raw := make([]interface{}, 4)
	raw[0] = "_error"
	raw[1] = transID
	raw[2] = props
	raw[3] = infos
	packet := Create.commandMessage(raw)
	return packet
}
