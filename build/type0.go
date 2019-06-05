package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// Type0 header
func Type0(timestamp, msglen, msgtypeid, msgstreamid int) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(timestamp, 3))
	buffer.Write(utils.WriteInt(msglen, 3))
	buffer.Write(utils.WriteInt(msgtypeid, 1))
	buffer.Write(utils.WriteInt(msgstreamid, 4))
	if timestamp >= 0xFFFFFF {
		buffer.Write(utils.WriteInt(timestamp, 4))
	}
	return buffer.Bytes()
}
