package main

import (
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
