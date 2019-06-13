package build

import (
	"SimpleRTMPServer/utils"
)

func (header) Basic(fmt, chunkid int) []byte {
	if chunkid <= 1 {
		return []byte{}
	}
	fmt = fmt & 3 << 6
	if chunkid < 64 {
		return []byte{byte(fmt | chunkid)}
	} else if chunkid < 65600 {
		return []byte{byte(fmt | 0), utils.WriteInt(chunkid-64, 1)[0]}
	}
	tmp := make([]byte, 3)
	tmp[0] = byte(fmt | 0)
	tmpval := utils.WriteInt(chunkid-64, 2)
	tmp[1] = tmpval[0]
	tmp[2] = tmpval[1]
	return tmp
}
