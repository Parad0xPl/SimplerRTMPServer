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

// ConnectionSettings Structure for stream data and settings
type ConnectionSettings struct {
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
}

// PacketData data
type PacketData struct {
	bytes []byte
}

// Packet wrapper for all packet related data
type Packet struct {
	ctx    *ConnectionSettings
	header *Header
	data   *PacketData
}

func initSettings() ConnectionSettings {
	return ConnectionSettings{
		ChunkSize:                   128,
		initTime:                    getTime(),
		ServerWindowAcknowledgement: 1024,
		PeerBandwidth:               128,
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
		settings.SizeRead += headers.Size
		log.Println("Headers", headers)
		if err != nil {
			log.Println(err)
			return
		}
		data := make([]byte, headers.MessageLength)
		_, err = c.Read(data)
		settings.SizeRead += headers.MessageLength
		if err != nil {
			log.Println("Error while reading body", err)
			return
		}

		packetData := PacketData{
			bytes: data,
		}

		packet := Packet{
			&settings,
			&headers,
			&packetData,
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
