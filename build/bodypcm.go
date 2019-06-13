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

func (pcmbody) AbortMessage(chunkid int) ([]byte, int) {
	return utils.WriteInt(chunkid, 4), 4
}

func (pcmbody) Acknowledgement(seqnumber int) ([]byte, int) {
	return utils.WriteInt(seqnumber, 4), 4
}

func (pcmbody) WindowAckSize(windowsize int) ([]byte, int) {
	return utils.WriteInt(windowsize, 4), 4
}

func (pcmbody) SetPeerBandwitdh(windowsize, limittype int) ([]byte, int) {
	buffer := new(bytes.Buffer)
	buffer.Write(utils.WriteInt(windowsize, 4))
	buffer.Write(utils.WriteInt(limittype, 1))
	ret := buffer.Bytes()
	return ret, len(ret)
}
