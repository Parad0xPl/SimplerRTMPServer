package connCTX

import (
	"SimpleRTMPServer/utils"
	"errors"
)

// StreamObject represent stream
type StreamObject struct{}

// ServerInstance Object which store server data
type ServerInstance struct {
	TopFreeID int
	NewConnID int
	DB        map[int]StreamObject
	Channels  map[string]*ChannelObject
	Published map[string]int
	Hash      *utils.HashGen
}

func (sm *ServerInstance) NewConn() int {
	defer func() {
		sm.NewConnID++
	}()
	return sm.NewConnID
}

func (sm *ServerInstance) CheckChannel(name string) bool {
	_, ok := sm.Channels[name]
	return ok
}

func (sm *ServerInstance) AddChannel(name, key string) {
	if sm.CheckChannel(name) == false {
		sm.Channels[name] = &ChannelObject{
			Name:       name,
			Key:        key,
			Subscribed: make(map[utils.Hash]*ConnContext),
		}
	}
}

func (sm *ServerInstance) CreateStream() int {
	tmp := sm.TopFreeID
	sm.TopFreeID++
	sm.DB[tmp] = StreamObject{}
	return tmp
}

func (sm *ServerInstance) CheckID(id int) bool {
	_, ok := sm.DB[id]
	return ok
}

func (sm *ServerInstance) DestroyStream(id int) {
	if sm.CheckID(id) {
		delete(sm.DB, id)
	}
}

func (sm *ServerInstance) CheckName(name string) bool {
	_, ok := sm.Published[name]
	return ok
}

func (sm *ServerInstance) Publish(id int, name string) error {
	if sm.CheckName(name) {
		return errors.New("Name reserved")
	}
	sm.Published[name] = id
	return nil
}

func (sm *ServerInstance) UnPublish(name string) {
	if !sm.CheckName(name) {
		return
	}
	delete(sm.Published, name)
}
