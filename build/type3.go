package build

import (
	"SimpleRTMPServer/utils"
)

// Type3 header
func (header) Type3(timestamp uint32) []byte {
	if timestamp > 0xffffff {
		return utils.WriteUintBE(uint(timestamp), 4)
	}
	return []byte{}
}
