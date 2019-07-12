package main

import (
	"SimpleRTMPServer/amf0"
	"fmt"
	"log"
	"net"
	"time"
)

func getTime() uint32 {
	return uint32(time.Now().UnixNano() / 1000)
}

// ReceivedPacket wrapper for all packet related data
type ReceivedPacket struct {
	ctx    *ConnContext
	header *Header
	data   []byte // Should work as long any handler won't try append to slice
}

func handler(conn net.Conn) {
	fmt.Printf("Connection started: %s\n", conn.RemoteAddr().String())
	ctx := initCTX(conn)
	defer ctx.Clear()

	// Handle handshake
	err := handshake(&ctx)
	if err != nil {
		log.Println(err)
		if conn.LocalAddr().Network() != "file" {
			return
		}
	}

netLoop:
	for {
		header, data, err := ctx.ReadPacket()
		if err != nil {
			log.Println("Error while reading Packet", err)
			return
		}

		packet := ReceivedPacket{
			&ctx,
			&header,
			data,
		}

		if header.ChunkID == 2 && header.StreamID == 0 {
			// Take effect when received
			err = handlePCM(packet)
			if err != nil {
				log.Println(err)
				break netLoop
			}
		} else {
			switch header.TypeID {
			case 4:
				// TODO handleUCM
				err = handleUCM(packet)
				if err != nil {
					log.Println(err)
					break netLoop
				}
			case 8:
				// handle Audio message
				if ctx.audioHandler != nil {
					ctx.audioHandler(data)
				}
			case 9:
				// TODO handle Video message
				if ctx.videoHandler != nil {
					ctx.videoHandler(data)
				}
			case 15:
				// TODO handle AMF3 data
			case 17:
				// TODO handle AMF3 command
			case 18:
				// TODO handle AMF0 data
				log.Println("AMF0 data received")
				parsedData := amf0.Read(data)
				i := 0
				parLen := len(parsedData)
				for i < parLen {
					item := parsedData[i]
					switch val := item.(type) {
					case string:
						if val == "onMetaData" {
							arr, ok := parsedData[i+1].(map[string]interface{})
							if !ok {
								log.Println("There is no metadata")
								break netLoop
							}
							i++
							fmt.Println("Metadata has been set")
							packet.ctx.channel.metadata = arr
						}
					}
					i++
				}
			case 20:
				err = handleAMF0cmd(packet)
				if err != nil {
					log.Println(err)
					break netLoop
				}
				// TODO handle AMF0 command

			}
		}
	}

}
