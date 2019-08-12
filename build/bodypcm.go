package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

type pcmbody struct{}

func (pcmbody) SetChunkSize(size int) ([]byte, int) {
	if size < 1 {
		size = 1
	}
	size = size & 0xefffff
	return utils.WriteIntBE(size, 4), 4
}

func (pcmbody) AbortMessage(chunkId int) ([]byte, int) {
	return utils.WriteIntBE(chunkId, 4), 4
}

func (pcmbody) Acknowledgement(seqNumber int) ([]byte, int) {
	return utils.WriteIntBE(seqNumber, 4), 4
}

func (pcmbody) WindowAckSize(windowSize int) ([]byte, int) {
	return utils.WriteIntBE(windowSize, 4), 4
}

func (pcmbody) SetPeerBandwidth(windowSize, limitType int) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteIntBE(windowSize, 4))
	buffer.Write(utils.WriteIntBE(limitType, 1))
	ret := buffer.Bytes()
	return ret, len(ret)
}
