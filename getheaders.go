package main

import (
	"errors"
	"net"
)

// Header options
type Header struct {
	ChunkID       int
	Format        int
	Timestamp     int
	MessageLength int
	TypeID        int
	StreamID      int
}

func headerEmpty() Header {
	return Header{}
}

func getExtandedTime(c net.Conn) (int, error) {
	tmp := make([]byte, 4)
	n, err := c.Read(tmp)
	if n != 4 || err != nil {
		return 0, errors.New("Problem with extanded time")
	}
	return readInt(tmp), nil
}

func getHeaders(c net.Conn, ctx *ConnectionSettings) (Header, error) {
	firstbyte := make([]byte, 1)
	var err error
	_, err = c.Read(firstbyte)
	if err != nil {
		return headerEmpty(), err
	}
	format := (int(firstbyte[0]) & 3 << 6) >> 6
	var chunkid int
	var tmp []byte
	if firstbyte[0] == 0 {
		tmp = make([]byte, 1)
		_, err = c.Read(tmp)
		if err != nil {
			return headerEmpty(), err
		}
		chunkid = readInt(tmp)
		chunkid += 64
	} else if firstbyte[0] == 1 {
		tmp = make([]byte, 2)
		_, err = c.Read(tmp)
		if err != nil {
			return headerEmpty(), err
		}
		chunkid = readInt(tmp)
		chunkid += 64
	} else if firstbyte[0] == 2 {
	} else {
		chunkid = int(firstbyte[0])
	}

	var timestamp, messlength, typeid, streamid int
	timestamp = -1
	messlength = -1
	typeid = -1
	streamid = -1
	if format == 0 {
		tmp = make([]byte, 11)
		_, err = c.Read(tmp)
		if err != nil {
			return headerEmpty(), err
		}
		timestamp = readInt(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtandedTime(c)
			if err != nil {
				return headerEmpty(), err
			}
			timestamp = t
		}
		messlength = readInt(tmp[3:6])
		typeid = int(tmp[7])
		streamid = readInt(tmp[8:])
	} else if format == 1 {
		tmp = make([]byte, 7)
		_, err = c.Read(tmp)
		if err != nil {
			return headerEmpty(), err
		}
		timestamp = readInt(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtandedTime(c)
			if err != nil {
				return headerEmpty(), err
			}
			timestamp = t
		}
		timestamp += ctx.lastHeader.Timestamp
		messlength = readInt(tmp[3:6])
		typeid = int(tmp[7])
		streamid = ctx.lastHeader.StreamID
	} else if format == 2 {
		tmp = make([]byte, 3)
		_, err = c.Read(tmp)
		if err != nil {
			return headerEmpty(), err
		}
		timestamp = readInt(tmp[0:3])
		if ^timestamp == 0 {
			t, err := getExtandedTime(c)
			if err != nil {
				return headerEmpty(), err
			}
			timestamp = t
		}
		timestamp += ctx.lastHeader.Timestamp
		streamid = ctx.lastHeader.StreamID
		messlength = ctx.lastHeader.MessageLength
		typeid = ctx.lastHeader.TypeID

	} else if format == 3 {
		timestamp = ctx.lastHeader.Timestamp
		streamid = ctx.lastHeader.StreamID
		messlength = ctx.lastHeader.MessageLength
		typeid = ctx.lastHeader.TypeID
	}

	lastHeader := Header{
		ChunkID:       chunkid,
		Format:        format,
		Timestamp:     timestamp,
		MessageLength: messlength,
		TypeID:        typeid,
		StreamID:      streamid,
	}
	return lastHeader, nil
}
