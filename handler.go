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
	data   []byte
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
		log.Println("Headers", headers)
		if err != nil {
			log.Println(err)
			return
		}
		data := make([]byte, headers.MessageLength)
		_, err = c.Read(data)
		if err != nil {
			log.Println("Error while reading body", err)
			return
		}

		packet := Packet{
			&settings,
			&headers,
			data,
		}
		if headers.ChunkID == 2 && headers.StreamID == 0 {
			// Take effect when received
			err = handlePCM(packet, c)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			switch headers.TypeID {
			case 4:
				// TODO handleUCM
				err = handleUCM(packet, c)
				if err != nil {
					log.Println(err)
					return
				}
			case 15:
				// TODO handle AMF3 data
			case 17:
				// TODO handle AMF3 command
			case 18:
				// TODO handle AMF0 data
			case 20:
				// TODO handle AMF0 command

			}
		}
	}

}
