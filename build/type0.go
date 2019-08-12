package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type0 header
func (header) Type0(timestamp uint32, msgLen, msgTypeID, msgStreamId int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteUintBE(uint(timestamp), 3))
	buffer.Write(utils.WriteIntBE(msgLen, 3))
	buffer.Write(utils.WriteIntBE(msgTypeID, 1))
	buffer.Write(utils.WriteIntLE(msgStreamId, 4))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteUintBE(uint(timestamp), 4))
	}
	return buffer.Bytes()
}
