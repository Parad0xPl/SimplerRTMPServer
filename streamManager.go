package main

import "errors"

// StreamObject represent stream
type StreamObject struct{}

// ChannelObject represent stream
type ChannelObject struct {
	name       string
	key        string
	subscribed map[string]*ConnContext
	metadata   map[string]interface{}
}

func (co *ChannelObject) subscribe(ctx *ConnContext) {
	// TODO name of subsriber generated from context
	co.subscribed["test"] = ctx
}

func (co *ChannelObject) unsubscribe(ctx *ConnContext) {
	// TODO name of subsriber generated from context
	delete(co.subscribed, "test")
}

func (co *ChannelObject) verify(key string) bool {
	return co.key == key
}

// StreamManager Object which manages streams
type StreamManager struct {
	topFreeID int
	db        map[int]StreamObject
	channels  map[string]*ChannelObject
	published map[string]int
}

func (sm *StreamManager) checkChannel(name string) bool {
	_, ok := sm.channels[name]
	return ok
}

func (sm *StreamManager) addChannel(name, key string) {
	if sm.checkChannel(name) == false {
		sm.channels[name] = &ChannelObject{
			name:       name,
			key:        key,
			subscribed: make(map[string]*ConnContext),
		}
	}
}

func (sm *StreamManager) createStream() int {
	tmp := sm.topFreeID
	sm.topFreeID++
	sm.db[tmp] = StreamObject{}
	return tmp
}

func (sm *StreamManager) checkid(id int) bool {
	_, ok := sm.db[id]
	return ok
}

func (sm *StreamManager) destroyStream(id int) {
	if sm.checkid(id) {
		delete(sm.db, id)
	}
}

func (sm *StreamManager) checkname(name string) bool {
	_, ok := sm.published[name]
	return ok
}

func (sm *StreamManager) publish(id int, name string) error {
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
		topFreeID: 10,
	}
	streamsMan.db = make(map[int]StreamObject)
	streamsMan.published = make(map[string]int)
	streamsMan.channels = make(map[string]*ChannelObject)

	// Temporary static channel
	streamsMan.addChannel("ksawk", "keyfortheapp")
}
