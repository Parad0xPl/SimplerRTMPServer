package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type2 header
func (header) Type2(timestamp int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(timestamp, 3))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteInt(timestamp, 4))
	}
	return buffer.Bytes()
}
