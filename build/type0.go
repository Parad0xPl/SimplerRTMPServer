package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type0 header
func (header) Type0(timestamp uint32, msgLen, msgTypeID, msgStreamId int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteUint(uint(timestamp), 3))
	buffer.Write(utils.WriteInt(msgLen, 3))
	buffer.Write(utils.WriteInt(msgTypeID, 1))
	buffer.Write(utils.WriteInt(msgStreamId, 4))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteUint(uint(timestamp), 4))
	}
	return buffer.Bytes()
}
