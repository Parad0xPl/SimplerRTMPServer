package main

import (
	"SimpleRTMPServer/connCTX"
	"SimpleRTMPServer/packet"
	"SimpleRTMPServer/utils"
	"fmt"
	"log"
	"net"
	"os"
)

func NewCTX(conn net.Conn) *connCTX.ConnContext {

	ctx := &connCTX.ConnContext{
		Index:                       serverInstance.NewConn(),
		Conn:                        conn,
		ChunkSize:                   128,
		InitTime:                    packet.RTMPTime(utils.GetTime()),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
		ChunkStreamID:               3,
		HashGen:                     serverInstance.Hash,
	}

	ctx.HeadersCache = packet.NewHeadersCache()

	if options.DumpInFnTemplate != "" &&
		options.DumpOutFnTemplate != "" {
		n := options.DumpFileCounter
		options.DumpFileCounter++

		readFilename := fmt.Sprintf("%s.%d", options.DumpInFnTemplate, n)
		writeFilename := fmt.Sprintf("%s.%d", options.DumpOutFnTemplate, n)

		fmt.Printf(
			"Opening dump files\nInput data: %s\nOutput data: %s\n",
			readFilename,
			writeFilename,
		)

		readFile, err := os.OpenFile(readFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Read dump file")
		} else {
			ctx.DumpFileForRead = readFile
		}

		writeFile, err := os.OpenFile(writeFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Write dump file")
		} else {
			ctx.DumpFileForWrite = writeFile
		}
	}

	return ctx
}
