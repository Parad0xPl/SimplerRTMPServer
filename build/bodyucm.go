package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

func (body) UCM(eventType int, data []byte) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteIntBE(eventType, 2))
	buffer.Write(data)
	return buffer.Bytes(), buffer.Len()
}
