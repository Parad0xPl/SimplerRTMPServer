package main

import (
	"SimpleRTMPServer/hash"
	"errors"
	"math/rand"
	"time"
)

// StreamObject represent stream
type StreamObject struct{}

// ChannelObject represent stream
// TODO: Clean() all closed connections
// TODO: Close() all subscribers
type ChannelObject struct {
	Name       string
	Key        string
	Subscribed map[hash.Type]*ConnContext
	Metadata   map[string]interface{}
}

func (co *ChannelObject) subscribe(ctx *ConnContext) {
	co.Subscribed[ctx.Hash()] = ctx
}

func (co *ChannelObject) isSubscribed(ctx *ConnContext) bool {
	_, ok := co.Subscribed[ctx.Hash()]
	return ok
}

func (co *ChannelObject) unsubscribe(ctx *ConnContext) {
	delete(co.Subscribed, ctx.Hash())
}

func (co *ChannelObject) verify(key string) bool {
	return co.Key == key
}

// ServerInstance Object which store server data
type ServerInstance struct {
	TopFreeID int
	NewConnID int
	DB        map[int]StreamObject
	Channels  map[string]*ChannelObject
	Published map[string]int
	Hash      hash.Gen
}

func (sm *ServerInstance) NewConn() int {
	defer func() {
		sm.NewConnID++
	}()
	return sm.NewConnID
}

func (sm *ServerInstance) checkChannel(name string) bool {
	_, ok := sm.Channels[name]
	return ok
}

func (sm *ServerInstance) addChannel(name, key string) {
	if sm.checkChannel(name) == false {
		sm.Channels[name] = &ChannelObject{
			Name:       name,
			Key:        key,
			Subscribed: make(map[hash.Type]*ConnContext),
		}
	}
}

func (sm *ServerInstance) createStream() int {
	tmp := sm.TopFreeID
	sm.TopFreeID++
	sm.DB[tmp] = StreamObject{}
	return tmp
}

func (sm *ServerInstance) checkID(id int) bool {
	_, ok := sm.DB[id]
	return ok
}

func (sm *ServerInstance) destroyStream(id int) {
	if sm.checkID(id) {
		delete(sm.DB, id)
	}
}

func (sm *ServerInstance) checkName(name string) bool {
	_, ok := sm.Published[name]
	return ok
}

func (sm *ServerInstance) publish(id int, name string) error {
	if sm.checkName(name) {
		return errors.New("Name reserved")
	}
	sm.Published[name] = id
	return nil
}

func (sm *ServerInstance) unpublish(name string) {
	if !sm.checkName(name) {
		return
	}
	delete(sm.Published, name)
}

var serverInstance ServerInstance

func initServerInstance() {

	rand.Seed(time.Now().UnixNano())

	serverInstance = ServerInstance{
		TopFreeID: 10,
		NewConnID: 1,
	}
	serverInstance.DB = make(map[int]StreamObject)
	serverInstance.Published = make(map[string]int)
	serverInstance.Channels = make(map[string]*ChannelObject)

	serverInstance.Hash = hash.InitGen()

	// Temporary static channel
	serverInstance.addChannel("ksawk", "keyfortheapp")
}
