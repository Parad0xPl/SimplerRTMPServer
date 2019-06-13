package build

import (
	"SimpleRTMPServer/utils"
)

// Type3 header
func (header) Type3(timestamp int) []byte {
	if timestamp > 0xffffff {
		return utils.WriteInt(timestamp, 4)
	}
	return []byte{}
}
