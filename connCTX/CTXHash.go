package connCTX

import (
	"SimpleRTMPServer/utils"
	"fmt"
)

func (ctx *ConnContext) Hash() utils.Hash {
	name := fmt.Sprintf(
		"[%d]%s",
		ctx.Index,
		ctx.Conn.RemoteAddr().String(),
	)
	return ctx.HashGen.String(name)
}
