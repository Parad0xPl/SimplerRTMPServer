package main

import (
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
	"strings"
)

// Header options
type Header struct {
	ChunkStreamID    int
	Format           int
	MessageTimestamp RTMPTime
	MessageLength    int
	MessageTypeID    int
	MessageStreamID  int
	Size             int
}

// Compare is headers have common values
func (h *Header) Compare(o Header) bool {
	if h.ChunkStreamID != o.ChunkStreamID {
		return false
	}
	if h.Format != o.Format {
		return false
	}
	if h.MessageLength != o.MessageLength {
		return false
	}
	if h.MessageStreamID != o.MessageStreamID {
		return false
	}
	if h.MessageTimestamp != o.MessageTimestamp {
		return false
	}
	if h.MessageTypeID != o.MessageTypeID {
		return false
	}
	return true
}

// Timestamp get timestamp
func (h *Header) Timestamp() uint32 {
	return h.MessageTimestamp.uint32()
}

// CopyFrom other header
func (h *Header) CopyFrom(o *Header) {
	h.ChunkStreamID = o.ChunkStreamID
	h.Format = o.Format
	h.MessageLength = o.MessageLength
	h.MessageStreamID = o.MessageStreamID
	h.MessageTimestamp = o.MessageTimestamp
	h.MessageTypeID = o.MessageTypeID
}

func (h Header) String() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Fmt:%d ", h.Format))
	str.WriteString(fmt.Sprintf("ChuID:%d ", h.ChunkStreamID))
	str.WriteString(fmt.Sprintf("Timestm:%d ", h.MessageTimestamp))
	str.WriteString(fmt.Sprintf("Msglen:%d ", h.MessageLength))
	str.WriteString(fmt.Sprintf("MessageTypeID:%d ", h.MessageTypeID))
	str.WriteString(fmt.Sprintf("MessageStrID:%d", h.MessageStreamID))
	return str.String()
}

func getExtendedTime(ctx *ConnContext) (uint, error) {
	tmp := make([]byte, 4)
	n, err := ctx.Read(tmp)
	if n != 4 || err != nil {
		return 0, errors.New("Problem with extended time.")
	}
	return utils.ReadUint(tmp), nil
}

// getHeader get header of RTMP chunk as specified in documentation(docs.pdf)
func getHeader(ctx *ConnContext) (Header, error) {
	// Read first byte with fmt and begin of chunk stream
	size := 0
	firstbyte := make([]byte, 1)
	var err error
	_, err = ctx.Read(firstbyte)
	size++
	if err != nil {
		return Header{}, err
	}
	// Splitting fmt of firstbyte
	var mask = 3 << 6
	format := (int(firstbyte[0]) & mask) >> 6
	mask = ^mask
	firstbyte[0] = byte(int(firstbyte[0]) & mask)
	var chunkid int
	var tmp []byte
	// Handling possible lengths
	switch firstbyte[0] {
	case 0:
		// 2 bytes long
		tmp = make([]byte, 1)
		_, err = ctx.Read(tmp)
		size++
		if err != nil {
			return Header{}, err
		}
		chunkid = utils.ReadInt(tmp)
		chunkid += 64
	case 1:
		// 3 bytes long
		tmp = make([]byte, 2)
		_, err = ctx.Read(tmp)
		size += 2
		if err != nil {
			return Header{}, err
		}
		chunkid = utils.ReadInt(tmp)
		chunkid += 64
	default:
		// 2 is reserved for low-level protocol control messages and commands
		chunkid = int(firstbyte[0])
	}

	var timestamp uint
	var messageLength, messageTypeID, messageStreamID int
	messageLength = -1
	messageTypeID = -1
	messageStreamID = -1
	switch format {
	case 0:
		// Type 0
		// 11 bytes long
		// Full data passed in header
		tmp = make([]byte, 11)
		_, err = ctx.Read(tmp)
		size += 11
		if err != nil {
			return Header{}, err
		}
		timestamp = utils.ReadUint(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtendedTime(ctx)
			size += 4
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		messageLength = utils.ReadInt(tmp[3:6])
		messageTypeID = int(tmp[6])
		messageStreamID = utils.ReadInt(tmp[7:])

	case 1:
		// Type 1
		// 7 bytes long
		// StreamID is skipped
		// Time delta instead of time
		tmp = make([]byte, 7)
		_, err = ctx.Read(tmp)
		size += 7
		if err != nil {
			return Header{}, err
		}
		timestamp = utils.ReadUint(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtendedTime(ctx)
			size += 4
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		timestamp += uint(ctx.LastHeaderReceived.MessageTimestamp)
		messageLength = utils.ReadInt(tmp[3:6])
		messageTypeID = int(tmp[6])
		messageStreamID = ctx.LastHeaderReceived.MessageStreamID

	case 2:
		// Type 2
		// 3 bytes long
		// Only time delta is given
		tmp = make([]byte, 3)
		_, err = ctx.Read(tmp)
		size += 3
		if err != nil {
			return Header{}, err
		}
		timestamp = utils.ReadUint(tmp[0:3])
		// TODO: fix condition
		if ^timestamp == 0 {
			size += 4
			t, err := getExtendedTime(ctx)
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		timestamp += uint(ctx.LastHeaderReceived.MessageTimestamp)
		messageStreamID = ctx.LastHeaderReceived.MessageStreamID
		messageLength = ctx.LastHeaderReceived.MessageLength
		messageTypeID = ctx.LastHeaderReceived.MessageTypeID

	case 3:
		// Type 3
		// Should use data from lastHeader
		timestamp = uint(ctx.LastHeaderReceived.MessageTimestamp)
		// TODO: fix condition
		if ^timestamp == 0 {
			size += 4
			t, err := getExtendedTime(ctx)
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		messageStreamID = ctx.LastHeaderReceived.MessageStreamID
		messageLength = ctx.LastHeaderReceived.MessageLength
		messageTypeID = ctx.LastHeaderReceived.MessageTypeID
	}

	lastHeader := Header{
		ChunkStreamID:    chunkid,
		Format:           format,
		MessageTimestamp: RTMPTime(timestamp),
		MessageLength:    messageLength,
		MessageTypeID:    messageTypeID,
		MessageStreamID:  messageStreamID,
		Size:             size,
	}
	ctx.LastHeaderReceived = lastHeader
	return lastHeader, nil
}
