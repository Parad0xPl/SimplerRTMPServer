package connCTX

import (
	"SimpleRTMPServer/utils"
)

// ChannelObject represent stream
// TODO: Clean() all closed connections
// TODO: Close() all subscribers
type ChannelObject struct {
	Name       string
	Key        string
	Subscribed map[utils.Hash]*ConnContext
	Metadata   map[string]interface{}
}

func (co *ChannelObject) Subscribe(ctx *ConnContext) {
	co.Subscribed[ctx.Hash()] = ctx
}

func (co *ChannelObject) IsSubscribed(ctx *ConnContext) bool {
	_, ok := co.Subscribed[ctx.Hash()]
	return ok
}

func (co *ChannelObject) Unsubscribe(ctx *ConnContext) {
	delete(co.Subscribed, ctx.Hash())
}

func (co *ChannelObject) Verify(key string) bool {
	return co.Key == key
}
