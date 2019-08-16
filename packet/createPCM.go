package packet

import "SimpleRTMPServer/build"

func (create) PCMSetChunkSize(size int) Prototype {
	head := &Header{
		MessageTypeID:   1,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.SetChunkSize(size)
	return Prototype{head, body}
}

func (create) PCMAbortMessage(chunkID int) Prototype {
	head := &Header{
		MessageTypeID:   2,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(chunkID)
	return Prototype{head, body}
}

func (create) PCMAcknowledgement(seqNumber int) Prototype {
	head := &Header{
		MessageTypeID:   3,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(seqNumber)
	return Prototype{head, body}
}

func (create) PCMWindowAckSize(winSize int) Prototype {
	head := &Header{
		MessageTypeID:   5,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.AbortMessage(winSize)
	return Prototype{head, body}
}

func (create) PCMSetPeerBandwidth(windowSize, limitType int) Prototype {
	head := &Header{
		MessageTypeID:   6,
		MessageStreamID: 0,
		ChunkStreamID:   2,
	}
	body, _ := build.Body.PCM.SetPeerBandwidth(windowSize, limitType)
	return Prototype{head, body}
}
