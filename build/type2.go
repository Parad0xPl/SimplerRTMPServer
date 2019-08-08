package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type2 header
func (header) Type2(timestamp uint32) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteUint(uint(timestamp), 3))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteUint(uint(timestamp), 4))
	}
	return buffer.Bytes()
}
