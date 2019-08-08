package main

import "errors"

// StreamObject represent stream
type StreamObject struct{}

// ChannelObject represent stream
// TODO: Clean() all closed connections
// TODO: Close() all subscribers
type ChannelObject struct {
	Name       string
	Key        string
	Subscribed map[string]*ConnContext
	Metadata   map[string]interface{}
}

func (co *ChannelObject) subscribe(ctx *ConnContext) {
	// TODO name of subsriber generated from context
	co.Subscribed["test"] = ctx
}

func (co *ChannelObject) unsubscribe(ctx *ConnContext) {
	// TODO name of subsriber generated from context
	delete(co.Subscribed, "test")
}

func (co *ChannelObject) verify(key string) bool {
	return co.Key == key
}

// StreamManager Object which manages streams
type StreamManager struct {
	TopFreeID int
	DB        map[int]StreamObject
	Channels  map[string]*ChannelObject
	Published map[string]int
}

func (sm *StreamManager) checkChannel(name string) bool {
	_, ok := sm.Channels[name]
	return ok
}

func (sm *StreamManager) addChannel(name, key string) {
	if sm.checkChannel(name) == false {
		sm.Channels[name] = &ChannelObject{
			Name:       name,
			Key:        key,
			Subscribed: make(map[string]*ConnContext),
		}
	}
}

func (sm *StreamManager) createStream() int {
	tmp := sm.TopFreeID
	sm.TopFreeID++
	sm.DB[tmp] = StreamObject{}
	return tmp
}

func (sm *StreamManager) checkid(id int) bool {
	_, ok := sm.DB[id]
	return ok
}

func (sm *StreamManager) destroyStream(id int) {
	if sm.checkid(id) {
		delete(sm.DB, id)
	}
}

func (sm *StreamManager) checkname(name string) bool {
	_, ok := sm.Published[name]
	return ok
}

func (sm *StreamManager) publish(id int, name string) error {
	if sm.checkname(name) {
		return errors.New("Name reserved")
	}
	sm.Published[name] = id
	return nil
}

func (sm *StreamManager) unpublish(name string) {
	if !sm.checkname(name) {
		return
	}
	delete(sm.Published, name)
}

var streamsMan StreamManager

func initStrMan() {
	streamsMan = StreamManager{
		TopFreeID: 10,
	}
	streamsMan.DB = make(map[int]StreamObject)
	streamsMan.Published = make(map[string]int)
	streamsMan.Channels = make(map[string]*ChannelObject)

	// Temporary static channel
	streamsMan.addChannel("ksawk", "keyfortheapp")
}
