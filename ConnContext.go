package main

import (
	"SimpleRTMPServer/utils"
	"fmt"
	"log"
	"net"
	"os"
)

type RawDataHandler func([]byte)
type RTMPTime uint32

func (t RTMPTime) uint32() uint32 {
	return uint32(t)
}

// ConnContext Structure for connection data and settings
type ConnContext struct {
	Index   int
	Conn    net.Conn
	Channel *ChannelObject

	DumpFileForRead  *os.File
	DumpFileForWrite *os.File

	AudioHandler RawDataHandler
	VideoHandler RawDataHandler

	ChunkSize     int
	InitTime      RTMPTime
	IsMetadataSet bool

	LastHeaderReceived *Header
	LastHeaderSend     *Header
	HeadersCache       *HeadersCache
	LastSendTimestamp  RTMPTime

	AmountRead    int
	LastAckAmount int
	AmountWrote   int

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

	if ctx.Channel != nil && ctx.Channel.IsSubscribed(ctx) {
		ctx.Channel.Unsubscribe(ctx)
	}
	ctx.Conn.Close()
}

// GetTime get current timestamp
func (ctx *ConnContext) GetTime() RTMPTime {
	return RTMPTime(utils.GetTime() - uint64(ctx.InitTime))
}

// Delta get delta of time
func (ctx *ConnContext) Delta(MessageTimestamp RTMPTime) uint32 {
	return (MessageTimestamp - ctx.LastSendTimestamp).uint32()
}

// Read Proxy for CTX.Conn.Read
func (ctx *ConnContext) Read(b []byte) (int, error) {
	n, err := ctx.Conn.Read(b)
	if ctx.DumpFileForRead != nil && err == nil {
		ctx.DumpFileForRead.Write(b[:n])
	}
	ctx.AmountRead += len(b)

	ctx.CheckAck()

	return n, err
}

// Write Proxy for CTX.Conn.Write
func (ctx *ConnContext) Write(b []byte) (int, error) {
	n, err := ctx.Conn.Write(b)
	if ctx.DumpFileForWrite != nil && err == nil {
		ctx.DumpFileForWrite.Write(b[:n])
	}
	ctx.AmountWrote += len(b)
	return n, err
}

func (ctx *ConnContext) CheckAck() {
	if ctx.ClientWindowAcknowledgement != 0 && ctx.AmountRead-ctx.LastAckAmount >= ctx.ClientWindowAcknowledgement {
		pkt := Create.PCMAcknowledgement(ctx.AmountRead - ctx.LastAckAmount)
		ctx.LastAckAmount = ctx.AmountRead
		ctx.SendPacket(pkt)
	}
}

// ReadPacket _
// TODO Refactor
func (ctx *ConnContext) ReadPacket() (Header, []byte, error) {
	header, err := getHeader(ctx)
	firstHeader := header
	if err != nil {
		return Header{}, []byte{}, err
	}
	dataToRead := header.MessageLength
	offset := 0
	chunkLen := utils.Min(dataToRead, ctx.ChunkSize)
	body := make([]byte, header.MessageLength)
	tmp := make([]byte, chunkLen)
	for {
		_, err = ctx.Read(tmp[:chunkLen])
		if err != nil {
			return Header{}, body, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = utils.Min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = getHeader(ctx)
		if err != nil {
			return Header{}, body, err
		}
	}
	return firstHeader, body, nil
}

func initCTX(conn net.Conn) ConnContext {

	ctx := ConnContext{
		Index:                       serverInstance.NewConn(),
		Conn:                        conn,
		ChunkSize:                   128,
		InitTime:                    RTMPTime(utils.GetTime()),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
		ChunkStreamID:               3,
	}

	ctx.HeadersCache = newHeadersCache()

	if options.DumpInFnTemplate != "" &&
		options.DumpOutFnTemplate != "" {
		n := options.DumpFileCounter
		options.DumpFileCounter++

		readFilename := fmt.Sprintf("%s.%d", options.DumpInFnTemplate, n)
		writeFilename := fmt.Sprintf("%s.%d", options.DumpOutFnTemplate, n)

		fmt.Printf(
			"Opening dump files\nInput data: %s\nOutput data: %s\n",
			readFilename,
			writeFilename,
		)

		readFile, err := os.OpenFile(readFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Read dump file")
		} else {
			ctx.DumpFileForRead = readFile
		}

		writeFile, err := os.OpenFile(writeFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Write dump file")
		} else {
			ctx.DumpFileForWrite = writeFile
		}
	}

	return ctx
}
