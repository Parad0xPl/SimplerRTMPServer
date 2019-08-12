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
	return utils.WriteInt(size, 4), 4
}

func (pcmbody) AbortMessage(chunkId int) ([]byte, int) {
	return utils.WriteInt(chunkId, 4), 4
}

func (pcmbody) Acknowledgement(seqNumber int) ([]byte, int) {
	return utils.WriteInt(seqNumber, 4), 4
}

func (pcmbody) WindowAckSize(windowSize int) ([]byte, int) {
	return utils.WriteInt(windowSize, 4), 4
}

func (pcmbody) SetPeerBandwidth(windowSize, limitType int) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(windowSize, 4))
	buffer.Write(utils.WriteInt(limitType, 1))
	ret := buffer.Bytes()
	return ret, len(ret)
}
