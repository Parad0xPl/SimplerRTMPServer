package connCTX

import (
	"SimpleRTMPServer/packet"
	"SimpleRTMPServer/utils"
	"net"
	"os"
)

type RawDataHandler func([]byte)

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
	InitTime      packet.RTMPTime
	IsMetadataSet bool

	LastHeaderReceived *packet.Header
	LastHeaderSend     *packet.Header
	HeadersCache       *packet.HeadersCache
	LastSendTimestamp  packet.RTMPTime

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

	HashGen *utils.HashGen
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
func (ctx *ConnContext) GetTime() packet.RTMPTime {
	return packet.RTMPTime(utils.GetTime() - uint64(ctx.InitTime))
}

// Delta get delta of time
func (ctx *ConnContext) Delta(MessageTimestamp packet.RTMPTime) uint32 {
	return (MessageTimestamp - ctx.LastSendTimestamp).Uint32()
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
		pkt := packet.Create.PCMAcknowledgement(ctx.AmountRead - ctx.LastAckAmount)
		ctx.LastAckAmount = ctx.AmountRead
		ctx.SendPacket(pkt)
	}
}

// ReadPacket _
// TODO Refactor
func (ctx *ConnContext) ReadPacket() (packet.Header, []byte, error) {
	header, err := ctx.getHeader()
	firstHeader := header
	if err != nil {
		return packet.Header{}, []byte{}, err
	}
	dataToRead := header.MessageLength
	offset := 0
	chunkLen := utils.Min(dataToRead, ctx.ChunkSize)
	body := make([]byte, header.MessageLength)
	tmp := make([]byte, chunkLen)
	for {
		_, err = ctx.Read(tmp[:chunkLen])
		if err != nil {
			return packet.Header{}, body, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = utils.Min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = ctx.getHeader()
		if err != nil {
			return packet.Header{}, body, err
		}
	}
	return firstHeader, body, nil
}
