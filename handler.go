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

// ConnContext Structure for stream data and settings
type ConnContext struct {
	ChunkSize                   int
	initTime                    uint32
	lastHeaderReceived          Header
	lastHeaderSended            *Header
	SizeRead                    int
	SizeWrote                   int
	Properties                  *map[string]interface{}
	ServerWindowAcknowledgement int
	ClientWindowAcknowledgement int
	PeerBandwidth               int
	PeerBandwidthType           int
}

// Packet wrapper for all packet related data
type Packet struct {
	ctx    *ConnContext
	header *Header
	data   []byte // Should work as long any handler won't try append to slice
}

func initCTX() ConnContext {
	return ConnContext{
		ChunkSize:                   128,
		initTime:                    getTime(),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
	}
}

func handler(c net.Conn) {
	defer c.Close()
	fmt.Printf("Connection started: %s\n", c.RemoteAddr().String())
	ctx := initCTX()

	// Handle handshake
	err := handshake(c)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		header, err := getHeader(c, &ctx)
		ctx.SizeRead += header.Size
		log.Println("Headers", header)
		if err != nil {
			log.Println(err)
			return
		}
		data := make([]byte, header.MessageLength)
		_, err = c.Read(data)
		ctx.SizeRead += header.MessageLength
		if err != nil {
			log.Println("Error while reading body", err)
			return
		}

		packet := Packet{
			&ctx,
			&header,
			data,
		}

		if header.ChunkID == 2 && header.StreamID == 0 {
			// Take effect when received
			err = handlePCM(packet, c)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			switch header.TypeID {
			case 4:
				// TODO handleUCM
				err = handleUCM(packet, c)
				if err != nil {
					log.Println(err)
					return
				}
			case 8:
				// TODO handle Audio message
			case 9:
				// TODO handle Video message
			case 15:
				// TODO handle AMF3 data
			case 17:
				// TODO handle AMF3 command
			case 18:
				// TODO handle AMF0 data
			case 20:
				handleAMF0cmd(packet, c)
				// TODO handle AMF0 command

			}
		}
	}

}
