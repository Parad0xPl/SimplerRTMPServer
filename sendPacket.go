package main

import (
	"SimpleRTMPServer/build"
	"bytes"
)

func checkType1(h, o Header) bool {
	if h.MessageStreamID != o.MessageStreamID {
		return false
	}
	return true
}

func checkType2(h, o Header) bool {
	if h.MessageLength != o.MessageLength {
		return false
	} else if h.MessageTypeID != o.MessageTypeID {
		return false
	}
	return true
}

func (ctx *ConnContext) sendChunk(pkt PacketPrototype) {
	header := pkt.Head
	body := pkt.Body
	buffer := new(bytes.Buffer)
	// TODO: header type base on some hash of difference
	// try to mark difference on exclusive bit
	if ctx.LastHeaderSend == nil {
		buffer.Write(build.Header.Basic(0, header.ChunkStreamID))
		buffer.Write(build.Header.Type0(
			header.Timestamp(),
			header.MessageLength,
			header.MessageTypeID,
			header.MessageStreamID,
		))
	} else if header.Compare(*ctx.LastHeaderSend) {
		buffer.Write(build.Header.Basic(3, header.ChunkStreamID))
		buffer.Write(build.Header.Type3(ctx.Delta(header.MessageTimestamp)))
	} else if checkType2(*ctx.LastHeaderSend, header) {
		buffer.Write(build.Header.Basic(2, header.ChunkStreamID))
		buffer.Write(build.Header.Type2(ctx.Delta(header.MessageTimestamp)))
	} else if checkType1(*ctx.LastHeaderSend, header) {
		buffer.Write(build.Header.Basic(1, header.ChunkStreamID))
		buffer.Write(build.Header.Type1(ctx.Delta(header.MessageTimestamp), header.MessageLength, header.MessageTypeID))
	} else {
		buffer.Write(build.Header.Basic(0, header.ChunkStreamID))
		buffer.Write(build.Header.Type0(header.Timestamp(), header.MessageLength, header.MessageTypeID, header.MessageStreamID))
	}
	ctx.LastHeaderSend = &header
	buffer.Write(body)
	ctx.Write(buffer.Bytes())
}

func (ctx *ConnContext) sendPacket(pkt PacketPrototype) {
	header := pkt.Head
	body := pkt.Body
	messLen := len(body)
	if header.ChunkStreamID == 0 {
		header.ChunkStreamID = ctx.ChunkStreamID
		defer func() {
			ctx.ChunkStreamID++
		}()
	}
	header.MessageLength = messLen
	header.MessageTimestamp = ctx.GetTime()
	for i := 0; i < messLen; i += ctx.ChunkSize {
		ctx.sendChunk(PacketPrototype{header, body[i:min(i+ctx.ChunkSize, messLen)]})
	}
	// TODO: support value overflow
	ctx.LastSendTimestamp = header.MessageTimestamp
}
