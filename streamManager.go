package main

import "errors"

// StreamObject represent stream
type StreamObject struct{}

// StreamManager Object which manages streams
type StreamManager struct {
	topFreeID uint
	db        map[uint]StreamObject
	published map[string]uint
}

func (sm *StreamManager) createStream() uint {
	tmp := sm.topFreeID
	sm.topFreeID++
	sm.db[tmp] = StreamObject{}
	return tmp
}

func (sm *StreamManager) checkid(id uint) bool {
	_, ok := sm.db[id]
	return ok
}

func (sm *StreamManager) destroyStream(id uint) {
	if sm.checkid(id) {
		delete(sm.db, id)
	}
}

func (sm *StreamManager) checkname(name string) bool {
	_, ok := sm.published[name]
	return ok
}

func (sm *StreamManager) publish(id uint, name string) error {
	if sm.checkname(name) {
		return errors.New("Name reserved")
	}
	sm.published[name] = id
	return nil
}

func (sm *StreamManager) unpublish(name string) {
	if !sm.checkname(name) {
		return
	}
	delete(sm.published, name)
}

var streamsMan StreamManager

func initStrMan() {
	streamsMan = StreamManager{
		topFreeID: 3,
	}
	streamsMan.db = make(map[uint]StreamObject)
	streamsMan.published = make(map[string]uint)
}
