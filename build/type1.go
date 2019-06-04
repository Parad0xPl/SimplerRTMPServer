package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type1 header
func Type1(timestamp, msglen, typeid int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(timestamp, 3))
	buffer.Write(utils.WriteInt(msglen, 3))
	buffer.Write(utils.WriteInt(typeid, 1))
	return buffer.Bytes()
}
