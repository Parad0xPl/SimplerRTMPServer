package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

func (body) UCM(eventtype int, data []byte) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(eventtype, 2))
	buffer.Write(data)
	return buffer.Bytes(), buffer.Len()
}
