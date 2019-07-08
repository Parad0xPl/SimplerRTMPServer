package main

import (
	"SimpleRTMPServer/build"
	"bytes"
	"net"
)

func checkType1(h, o Header) bool {
	if h.StreamID != o.StreamID {
		return false
	}
	return true
}

func checkType2(h, o Header) bool {
	if h.MessageLength != o.MessageLength {
		return false
	} else if h.TypeID != o.TypeID {
		return false
	}
	return true
}

func sendPacket(c net.Conn, ctx *ConnContext, pkt PacketProt) {
	header := pkt.head
	body := pkt.body
	buffer := new(bytes.Buffer)
	messLen := len(body)
	header.MessageLength = messLen
	if ctx.lastHeaderSended == nil {
		buffer.Write(build.Header.Basic(0, header.ChunkID))
		buffer.Write(build.Header.Type0(header.Timestamp, header.MessageLength, header.TypeID, header.StreamID))
	} else if header.Compare(*ctx.lastHeaderSended) {
		buffer.Write(build.Header.Basic(3, header.ChunkID))
		buffer.Write(build.Header.Type3(header.Timestamp))
	} else if checkType2(*ctx.lastHeaderSended, header) {
		buffer.Write(build.Header.Basic(2, header.ChunkID))
		buffer.Write(build.Header.Type2(header.Timestamp))
	} else if checkType1(*ctx.lastHeaderSended, header) {
		buffer.Write(build.Header.Basic(1, header.ChunkID))
		buffer.Write(build.Header.Type1(header.Timestamp, header.MessageLength, header.TypeID))
	} else {
		buffer.Write(build.Header.Basic(0, header.ChunkID))
		buffer.Write(build.Header.Type0(header.Timestamp, header.MessageLength, header.TypeID, header.StreamID))
	}
	ctx.lastHeaderSended = &header
	buffer.Write(body)
	ctx.SizeWrote += buffer.Len()
	c.Write(buffer.Bytes())
}
