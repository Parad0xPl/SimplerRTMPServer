package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type1 header
func (header) Type1(timestamp uint32, msgLen, typeId int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteUintBE(uint(timestamp), 3))
	buffer.Write(utils.WriteIntBE(msgLen, 3))
	buffer.Write(utils.WriteIntBE(typeId, 1))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteUintBE(uint(timestamp), 4))
	}
	return buffer.Bytes()
}
