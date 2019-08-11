package main

import (
	"SimpleRTMPServer/utils"
	"errors"
)

// HeaderReceiver provide nice way of reading header from ctx
// handles reading data from last header with same chunk stream ID
type HeaderReceiver struct {
	ctx  *ConnContext
	size int

	isBasicHeaderRead   bool
	isMessageHeaderRead bool

	lastHeader *Header

	chunkID int
	format  int

	timestamp       uint
	messageLength   int
	messageTypeID   int
	messageStreamID int
}

// Header return finished header
func (builder *HeaderReceiver) Header() Header {
	return Header{
		ChunkStreamID:    builder.chunkID,
		Format:           builder.format,
		MessageTimestamp: RTMPTime(builder.timestamp),
		MessageLength:    builder.messageLength,
		MessageTypeID:    builder.messageTypeID,
		MessageStreamID:  builder.messageStreamID,
		Size:             builder.size,
	}
}

// getBasicHeader read format and chunk stream id
// values are places in structure
func (builder *HeaderReceiver) getBasicHeader() error {
	if builder.isBasicHeaderRead == true {
		return errors.New("Basic header already read.")
	}
	var err error

	firstByte := make([]byte, 1)
	_, err = builder.ctx.Read(firstByte)
	builder.size++
	if err != nil {
		return err
	}

	// Splitting fmt of firstByte
	var mask = 3 << 6
	format := (int(firstByte[0]) & mask) >> 6
	mask = ^mask
	firstByte[0] = byte(int(firstByte[0]) & mask)

	var chunkID int
	var tmp []byte

	// Handling possible lengths
	switch firstByte[0] {
	case 0:
		// 2 bytes long
		tmp = make([]byte, 1)
		_, err = builder.ctx.Read(tmp)
		builder.size++
		if err != nil {
			return err
		}
		chunkID = utils.ReadInt(tmp)
		chunkID += 64
	case 1:
		// 3 bytes long
		tmp = make([]byte, 2)
		_, err = builder.ctx.Read(tmp)
		builder.size += 2
		if err != nil {
			return err
		}
		chunkID = utils.ReadInt(tmp)
		chunkID += 64
	default:
		chunkID = int(firstByte[0])
	}
	builder.format = format
	builder.chunkID = chunkID
	builder.isBasicHeaderRead = true
	return nil
}

// getExtendedTime read extended time and replace old value
func (builder *HeaderReceiver) getExtendedTime() error {
	tmp := make([]byte, 4)
	n, err := builder.ctx.Read(tmp)
	if n != 4 || err != nil {
		return errors.New("Problem with extended time.")
	}
	builder.size += 4
	builder.timestamp = utils.ReadUint(tmp)
	return nil
}

// type0 read specific format
func (builder *HeaderReceiver) type0() error {
	// Type 0
	// 11 bytes long
	// Full data passed in header
	tmp := make([]byte, 11)
	_, err := builder.ctx.Read(tmp)
	builder.size += 11
	if err != nil {
		return err
	}
	builder.timestamp = utils.ReadUint(tmp[0:3])
	if builder.timestamp == 0xFFFFFF {
		err := builder.getExtendedTime()
		if err != nil {
			return err
		}
	}
	builder.messageLength = utils.ReadInt(tmp[3:6])
	builder.messageTypeID = int(tmp[6])
	builder.messageStreamID = utils.ReadInt(tmp[7:])
	return nil
}

// type1 read specific format
func (builder *HeaderReceiver) type1() error {
	// Type 1
	// 7 bytes long
	// StreamID is skipped
	// Time delta instead of time
	tmp := make([]byte, 7)
	_, err := builder.ctx.Read(tmp)
	builder.size += 7
	if err != nil {
		return err
	}
	builder.timestamp = utils.ReadUint(tmp[0:3])
	if builder.timestamp == 0xFFFFFF {
		err := builder.getExtendedTime()
		if err != nil {
			return err
		}
	}
	builder.timestamp += uint(builder.lastHeader.MessageTimestamp)
	builder.messageLength = utils.ReadInt(tmp[3:6])
	builder.messageTypeID = int(tmp[6])
	builder.messageStreamID = builder.lastHeader.MessageStreamID
	return nil
}

// type2 read specific format
func (builder *HeaderReceiver) type2() error {
	// Type 2
	// 3 bytes long
	// Only time delta is given
	tmp := make([]byte, 3)
	_, err := builder.ctx.Read(tmp)
	builder.size += 3
	if err != nil {
		return err
	}
	builder.timestamp = utils.ReadUint(tmp[0:3])
	if builder.timestamp == 0xFFFFFF {
		err := builder.getExtendedTime()
		if err != nil {
			return err
		}
	}
	builder.timestamp += uint(builder.lastHeader.MessageTimestamp)
	builder.messageStreamID = builder.lastHeader.MessageStreamID
	builder.messageLength = builder.lastHeader.MessageLength
	builder.messageTypeID = builder.lastHeader.MessageTypeID
	return nil
}

// type3 read specific format
func (builder *HeaderReceiver) type3() error {
	// Type 3
	// Should use data from builder.lastHeader
	builder.timestamp = uint(builder.lastHeader.MessageTimestamp)
	if builder.timestamp == 0xFFFFFF {
		builder.size += 4
		err := builder.getExtendedTime()
		if err != nil {
			return err
		}
	}
	builder.messageStreamID = builder.lastHeader.MessageStreamID
	builder.messageLength = builder.lastHeader.MessageLength
	builder.messageTypeID = builder.lastHeader.MessageTypeID
	return nil
}

// messageHeader read message header
func (builder *HeaderReceiver) messageHeader() error {
	if builder.isMessageHeaderRead == true {
		return errors.New("Message header already read.")
	}
	var err error
	switch builder.format {
	case 0:
		err = builder.type0()
		if err != nil {
			return err
		}
	case 1:
		err = builder.type1()
		if err != nil {
			return err
		}

	case 2:
		err = builder.type2()
		if err != nil {
			return err
		}

	case 3:
		err = builder.type3()
		if err != nil {
			return err
		}

	}
	builder.isMessageHeaderRead = true
	return nil
}

// Get get header of RTMP chunk as specified in documentation(docs.pdf)
func (builder *HeaderReceiver) Get() error {
	// Read first byte with fmt and begin of chunk stream
	err := builder.getBasicHeader()
	if err != nil {
		return err
	}

	builder.messageLength = -1
	builder.messageTypeID = -1
	builder.messageStreamID = -1

	builder.lastHeader = builder.ctx.HeadersCache.Get(builder.chunkID)
	if builder.lastHeader == nil {
		builder.lastHeader = builder.ctx.LastHeaderReceived
	}

	err = builder.messageHeader()
	if err != nil {
		return err
	}

	return nil
}

// getHeader read Header from ctx
func getHeader(ctx *ConnContext) (Header, error) {
	builder := HeaderReceiver{
		ctx: ctx,
	}
	err := builder.Get()
	if err != nil {
		return Header{}, err
	}

	header := builder.Header()
	ctx.HeadersCache.Insert(header.ChunkStreamID, &header)
	ctx.LastHeaderReceived = &header

	return header, nil
}
