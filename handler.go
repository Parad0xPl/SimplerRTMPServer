package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func getTime() uint32 {
	return uint32(time.Now().UnixNano() / 1000)
}

func readInt(b []byte) int {
	tmp := 0
	for i := 0; i < len(b); i++ {
		tmp = tmp<<8 | int(b[i])
	}
	return tmp
}

// ConnectionSettings Structure for stream data and settings
type ConnectionSettings struct {
	ChunkSize  int
	initTime   uint32
	lastHeader Header
}

// Packet wrapper for all packet related data
type Packet struct {
	ctx    *ConnectionSettings
	header *Header
}

func initSettings() ConnectionSettings {
	return ConnectionSettings{
		ChunkSize: 128,
		initTime:  getTime(),
	}
}

func handler(c net.Conn) {
	defer c.Close()
	fmt.Printf("Connection started: %s\n", c.RemoteAddr().String())
	settings := initSettings()

	err := handshake(c)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		headers, err := getHeaders(c, &settings)
		log.Println(headers)
		if err != nil {
			log.Println(err)
			return
		}
		packet := Packet{
			&settings,
			&headers,
		}
		if headers.ChunkID == 2 && headers.StreamID == 0 {
			// Take effect when received
			err = handlePCM(packet, c)
			if err != nil {
				log.Println(err)
				return
			}
		} else if headers.TypeID == 4 {
			// TODO handleUCM
			err = handleUCM(packet, c)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}
