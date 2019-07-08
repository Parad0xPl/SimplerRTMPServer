package main

// StreamObject represent stream
type StreamObject struct{}

// StreamManager Object which manages streams
type StreamManager struct {
	topFreeID uint
	db        map[uint]StreamObject
}

func (sm *StreamManager) createStream() uint {
	tmp := sm.topFreeID
	sm.topFreeID++
	sm.db[tmp] = StreamObject{}
	return tmp
}

func (sm *StreamManager) check(id uint) bool {
	_, ok := sm.db[id]
	return ok
}

func (sm *StreamManager) destroyStream(id uint) {
	if sm.check(id) {
		delete(sm.db, id)
	}
}

var streamsMan StreamManager

func initStrMan() {
	streamsMan = StreamManager{
		topFreeID: 3,
	}
}
