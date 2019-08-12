package main

import (
	"SimpleRTMPServer/hash"
	"fmt"
)

func (ctx *ConnContext) Hash() hash.Type {
	name := fmt.Sprintf(
		"[%d]%s",
		ctx.Index,
		ctx.Conn.RemoteAddr().String(),
	)
	return serverInstance.Hash.String(name)
}
