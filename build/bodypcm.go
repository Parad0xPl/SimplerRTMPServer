package build

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

type pcmbody struct{}

func (pcmbody) pcmSetChunkSize(size int) ([]byte, int) {
	if size < 1 {
		size = 1
	}
	size = size & 0xefffff
	return utils.WriteInt(size, 4), 4
}

func (pcmbody) pcmAbortMessage(chunkid int) ([]byte, int) {
	return utils.WriteInt(chunkid, 4), 4
}

func (pcmbody) pcmAcknowledgement(seqnumber int) ([]byte, int) {
	return utils.WriteInt(seqnumber, 4), 4
}

func (pcmbody) pcmWindowAckSize(windowsize int) ([]byte, int) {
	return utils.WriteInt(windowsize, 4), 4
}

func (pcmbody) pcmSetPeerBandwitdh(windowsize, limittype int) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(windowsize, 4))
	buffer.Write(utils.WriteInt(limittype, 1))
	ret := buffer.Bytes()
	return ret, len(ret)
}
