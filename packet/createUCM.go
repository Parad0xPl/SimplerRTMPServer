package packet

import (
	"SimpleRTMPServer/build"
	"SimpleRTMPServer/utils"
)

func (create) UCMStreamBegin(streamID int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(0, utils.WriteIntBE(streamID, 4))
	return Prototype{head, body}
}

func (create) UCMStreamEOF(streamID int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(1, utils.WriteIntBE(streamID, 4))
	return Prototype{head, body}
}

func (create) UCMStreamDry(streamID int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(2, utils.WriteIntBE(streamID, 4))
	return Prototype{head, body}
}

func (create) UCMSetBufferLength(streamID, buffLen int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	eventData := utils.Concat(utils.WriteIntBE(streamID, 4), utils.WriteIntBE(buffLen, 4))
	body, _ := build.Body.UCM(3, eventData)
	return Prototype{head, body}
}

func (create) UCMStreamIsRecorded(streamID int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(4, utils.WriteIntBE(streamID, 4))
	return Prototype{head, body}
}

func (create) UCMPingRequest(timestamp int) Prototype {
	head := &Header{
		MessageTypeID:   4,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.UCM(6, utils.WriteIntBE(timestamp, 4))
	return Prototype{head, body}
}

func (create) UCMPingResponse(timestamp interface{}) Prototype {
	head := &Header{
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
	return Prototype{head, body}
}
