package connCTX

import (
	"SimpleRTMPServer/packet"
)

// getHeader read Header from ctx
func (ctx *ConnContext) getHeader() (packet.Header, error) {
	builder := HeaderReceiver{
		CTX: ctx,
	}
	err := builder.Get()
	if err != nil {
		return packet.Header{}, err
	}

	header := builder.Header()
	ctx.HeadersCache.Insert(header.ChunkStreamID, &header)
	ctx.LastHeaderReceived = &header

	return header, nil
}
