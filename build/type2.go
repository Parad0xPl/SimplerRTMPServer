package build

import (
	"SimpleRTMPServer/utils"
)

// Type2 header
func Type2(timestamp int) []byte {
	return utils.WriteInt(timestamp, 3)
}
