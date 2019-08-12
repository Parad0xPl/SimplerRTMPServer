package build

import (
	"SimpleRTMPServer/utils"
)

func (header) Basic(fmt, chunkID int) []byte {
	if chunkID <= 1 {
		return []byte{}
	}
	fmt = fmt & 3 << 6
	if chunkID < 64 {
		return []byte{byte(fmt | chunkID)}
	} else if chunkID < 65600 {
		return []byte{byte(fmt | 0), utils.WriteIntBE(chunkID-64, 1)[0]}
	}
	tmp := make([]byte, 3)
	tmp[0] = byte(fmt | 0)
	tmpVal := utils.WriteIntBE(chunkID-64, 2)
	tmp[1] = tmpVal[0]
	tmp[2] = tmpVal[1]
	return tmp
}
