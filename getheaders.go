package main

import (
	"SimpleRTMPServer/utils"
	"errors"
	"fmt"
	"strings"
)

// Header options
type Header struct {
	ChunkID       int
	Format        int
	Timestamp     int
	MessageLength int
	TypeID        int
	StreamID      int
	Size          int
}

// Compare is headers have common values
func (h *Header) Compare(o Header) bool {
	if h.ChunkID != o.ChunkID {
		return false
	}
	if h.Format != o.Format {
		return false
	}
	if h.MessageLength != o.MessageLength {
		return false
	}
	if h.StreamID != o.StreamID {
		return false
	}
	if h.Timestamp != o.Timestamp {
		return false
	}
	if h.TypeID != o.TypeID {
		return false
	}
	return true
}

// CopyFrom other header
func (h *Header) CopyFrom(o *Header) {
	h.ChunkID = o.ChunkID
	h.Format = o.Format
	h.MessageLength = o.MessageLength
	h.StreamID = o.StreamID
	h.Timestamp = o.Timestamp
	h.TypeID = o.TypeID
}

func (h Header) String() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Fmt:%d ", h.Format))
	str.WriteString(fmt.Sprintf("ChuID:%d ", h.ChunkID))
	str.WriteString(fmt.Sprintf("Timestm:%d ", h.Timestamp))
	str.WriteString(fmt.Sprintf("Msglen:%d ", h.MessageLength))
	str.WriteString(fmt.Sprintf("TypeID:%d ", h.TypeID))
	str.WriteString(fmt.Sprintf("StrID:%d", h.StreamID))
	return str.String()
}

func getExtandedTime(ctx *ConnContext) (int, error) {
	tmp := make([]byte, 4)
	n, err := ctx.Read(tmp)
	if n != 4 || err != nil {
		return 0, errors.New("Problem with extanded time")
	}
	return utils.ReadInt(tmp), nil
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

	var timestamp, messlength, typeid, streamid int
	timestamp = -1
	messlength = -1
	typeid = -1
	streamid = -1
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
		timestamp = utils.ReadInt(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtandedTime(ctx)
			size += 4
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		messlength = utils.ReadInt(tmp[3:6])
		typeid = int(tmp[6])
		streamid = utils.ReadInt(tmp[7:])
	case 1:
		// Type 1
		// 7 bytes long
		// Streamid is sipped
		// Timedelta instead of time
		tmp = make([]byte, 7)
		_, err = ctx.Read(tmp)
		size += 7
		if err != nil {
			return Header{}, err
		}
		timestamp = utils.ReadInt(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtandedTime(ctx)
			size += 4
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		timestamp += ctx.lastHeaderReceived.Timestamp
		messlength = utils.ReadInt(tmp[3:6])
		typeid = int(tmp[6])
		streamid = ctx.lastHeaderReceived.StreamID
	case 2:
		// Type 2
		// 3 bytes long
		// Only timedelta is given
		tmp = make([]byte, 3)
		_, err = ctx.Read(tmp)
		size += 3
		if err != nil {
			return Header{}, err
		}
		timestamp = utils.ReadInt(tmp[0:3])
		if ^timestamp == 0 {
			size += 4
			t, err := getExtandedTime(ctx)
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		timestamp += ctx.lastHeaderReceived.Timestamp
		streamid = ctx.lastHeaderReceived.StreamID
		messlength = ctx.lastHeaderReceived.MessageLength
		typeid = ctx.lastHeaderReceived.TypeID

	case 3:
		// Type 3
		// Should use data from lastHeader
		timestamp = ctx.lastHeaderReceived.Timestamp
		if ^timestamp == 0 {
			size += 4
			t, err := getExtandedTime(ctx)
			if err != nil {
				return Header{}, err
			}
			timestamp = t
		}
		streamid = ctx.lastHeaderReceived.StreamID
		messlength = ctx.lastHeaderReceived.MessageLength
		typeid = ctx.lastHeaderReceived.TypeID
	}

	lastHeader := Header{
		ChunkID:       chunkid,
		Format:        format,
		Timestamp:     timestamp,
		MessageLength: messlength,
		TypeID:        typeid,
		StreamID:      streamid,
		Size:          size,
	}
	ctx.lastHeaderReceived = lastHeader
	return lastHeader, nil
}
