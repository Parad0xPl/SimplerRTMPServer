package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type1 header
func (header) Type1(timestamp uint32, msgLen, typeId int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteUint(uint(timestamp), 3))
	buffer.Write(utils.WriteInt(msgLen, 3))
	buffer.Write(utils.WriteInt(typeId, 1))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteUint(uint(timestamp), 4))
	}
	return buffer.Bytes()
}
