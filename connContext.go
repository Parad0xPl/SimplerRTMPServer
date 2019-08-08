package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type RawDataHandler func([]byte)

// ConnContext Structure for connection data and settings
type ConnContext struct {
	Conn    net.Conn
	Channel *ChannelObject

	DumpFileForRead  *os.File
	DumpFileForWrite *os.File

	AudioHandler RawDataHandler
	VideoHandler RawDataHandler

	ChunkSize     int
	InitTime      uint32
	IsMetadataSet bool

	LastHeaderReceived Header
	LastHeaderSend     *Header

	AmountRead  int
	AmountWrote int

	Properties *map[string]interface{}

	ServerWindowAcknowledgement int
	ClientWindowAcknowledgement int
	PeerBandwidth               int
	PeerBandwidthType           int

	ChunkStreamID int
	NetStreamID   int
}

// Clear context before destructing
func (ctx *ConnContext) Clear() {
	if ctx.DumpFileForRead != nil {
		ctx.DumpFileForRead.Close()
	}
	if ctx.DumpFileForWrite != nil {
		ctx.DumpFileForWrite.Close()
	}
	if ctx.IsMetadataSet {
		ctx.Channel.Metadata = nil
	}
	ctx.Conn.Close()
}

// Read Proxy for CTX.Conn.Read
func (ctx ConnContext) Read(b []byte) (int, error) {
	n, err := ctx.Conn.Read(b)
	if ctx.DumpFileForRead != nil {
		ctx.DumpFileForRead.Write(b[:n])
	}
	ctx.AmountRead += len(b)
	return n, err
}

// Write Proxy for CTX.Conn.Write
func (ctx ConnContext) Write(b []byte) (int, error) {
	n, err := ctx.Conn.Write(b)
	if ctx.DumpFileForWrite != nil {
		ctx.DumpFileForWrite.Write(b[:n])
	}
	ctx.AmountWrote += len(b)
	return n, err
}

// ReadPacket _
// TODO Refactor
func (ctx *ConnContext) ReadPacket() (Header, []byte, error) {
	header, err := getHeader(ctx)
	firstheader := header
	if err != nil {
		return Header{}, []byte{}, err
	}
	dataToRead := header.MessageLength
	offset := 0
	chunkLen := min(dataToRead, ctx.ChunkSize)
	body := make([]byte, header.MessageLength)
	tmp := make([]byte, chunkLen)
	for {
		_, err = ctx.Read(tmp[:chunkLen])
		if err != nil {
			return Header{}, body, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = getHeader(ctx)
		if err != nil {
			return Header{}, body, err
		}
	}
	return firstheader, body, nil
}

func initCTX(conn net.Conn) ConnContext {

	ctx := ConnContext{
		Conn:                        conn,
		ChunkSize:                   128,
		InitTime:                    uint32(getTime()),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
		ChunkStreamID:               3,
	}

	if options.DumpInFnTemplate != "" &&
		options.DumpOutFnTemplate != "" {
		n := options.DumpFileCounter
		options.DumpFileCounter++

		readfilename := fmt.Sprintf("%s.%d", options.DumpInFnTemplate, n)
		writefilename := fmt.Sprintf("%s.%d", options.DumpOutFnTemplate, n)

		fmt.Printf(
			"Opening dump files\nInput data: %s\nOutput data: %s\n",
			readfilename,
			writefilename,
		)

		readfile, err := os.OpenFile(readfilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Read dump file")
		} else {
			ctx.DumpFileForRead = readfile
		}

		writefile, err := os.OpenFile(writefilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Write dump file")
		} else {
			ctx.DumpFileForWrite = writefile
		}
	}

	return ctx
}
