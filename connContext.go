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

type rawDataHandler func([]byte)

// ConnContext Structure for connection data and settings
type ConnContext struct {
	conn    net.Conn
	channel *ChannelObject

	readDumpFile  *os.File
	writeDumpFile *os.File

	audioHandler rawDataHandler
	videoHandler rawDataHandler

	ChunkSize int
	initTime  uint32

	lastHeaderReceived Header
	lastHeaderSended   *Header

	SizeRead  int
	SizeWrote int

	Properties *map[string]interface{}

	ServerWindowAcknowledgement int
	ClientWindowAcknowledgement int
	PeerBandwidth               int
	PeerBandwidthType           int

	StreamID int
}

// Clear context before destructing
func (ctx *ConnContext) Clear() {
	if ctx.readDumpFile != nil {
		ctx.readDumpFile.Close()
	}
	if ctx.writeDumpFile != nil {
		ctx.writeDumpFile.Close()
	}
	ctx.conn.Close()
}

// Read Proxy for ctx.conn.Read
func (ctx ConnContext) Read(b []byte) (int, error) {
	n, err := ctx.conn.Read(b)
	if ctx.readDumpFile != nil {
		ctx.readDumpFile.Write(b[:n])
	}
	return n, err
}

// Write Proxy for ctx.conn.Write
func (ctx ConnContext) Write(b []byte) (int, error) {
	n, err := ctx.conn.Write(b)
	if ctx.writeDumpFile != nil {
		ctx.writeDumpFile.Write(b[:n])
	}
	return n, err
}

// ReadPacket _
func (ctx *ConnContext) ReadPacket() (Header, []byte, error) {
	header, err := getHeader(ctx)
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
			return Header{}, []byte{}, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = getHeader(ctx)
		if err != nil {
			return Header{}, []byte{}, err
		}
	}
	return header, body, nil
}

func initCTX(conn net.Conn) ConnContext {

	ctx := ConnContext{
		conn:                        conn,
		ChunkSize:                   128,
		initTime:                    getTime(),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
	}

	if options.dumpfilein != "" &&
		options.dumpfileout != "" {
		n := options.dumpfilecounter
		options.dumpfilecounter++

		readfilename := fmt.Sprintf("%s.%d", options.dumpfilein, n)
		writefilename := fmt.Sprintf("%s.%d", options.dumpfileout, n)

		fmt.Printf(
			"Opening dump files\nInput data: %s\nOutput data: %s\n",
			readfilename,
			writefilename,
		)

		readfile, err := os.OpenFile(readfilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Read dump file")
		} else {
			ctx.readDumpFile = readfile
		}

		writefile, err := os.OpenFile(writefilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Write dump file")
		} else {
			ctx.writeDumpFile = writefile
		}
	}

	return ctx
}
