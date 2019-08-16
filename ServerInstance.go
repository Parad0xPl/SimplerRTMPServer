package main

import (
	"SimpleRTMPServer/connCTX"
	"SimpleRTMPServer/utils"
	"math/rand"
	"time"
)

var serverInstance *connCTX.ServerInstance

func initServerInstance() {

	rand.Seed(time.Now().UnixNano())

	serverInstance = &connCTX.ServerInstance{
		TopFreeID: 10,
		NewConnID: 1,
	}
	serverInstance.DB = make(map[int]connCTX.StreamObject)
	serverInstance.Published = make(map[string]int)
	serverInstance.Channels = make(map[string]*connCTX.ChannelObject)

	serverInstance.Hash = utils.NewGen()

	// Temporary static channel
	serverInstance.AddChannel("ksawk", "keyfortheapp")
}
