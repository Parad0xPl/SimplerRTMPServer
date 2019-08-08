package main

import (
	"SimpleRTMPServer/amf0"
	"fmt"
	"log"
	"net"
)

// ReceivedPacket wrapper for all packet related data
type ReceivedPacket struct {
	CTX    *ConnContext
	Header *Header
	Data   []byte // Should work as long any handler won't try append to slice
}

func handlePacket(packet ReceivedPacket) error {
	var err error
	if packet.Header.ChunkStreamID == 2 && packet.Header.MessageStreamID == 0 {
		// Take effect when received
		err = handlePCM(packet)
		return err
	} else {
		switch packet.Header.MessageTypeID {
		case 8:
			if packet.CTX.AudioHandler != nil {
				log.Println("Audio data received. Len:", len(packet.Data))
				packet.CTX.AudioHandler(packet.Data)
			}
			return nil
		case 9:
			if packet.CTX.VideoHandler != nil {
				log.Println("Video data received. Len:", len(packet.Data))
				packet.CTX.VideoHandler(packet.Data)
			}
			return nil
		case 15:
			// TODO handle AMF3 data
		case 17:
			// TODO handle AMF3 command
		case 18:
			log.Println("AMF0 data received")
			parsedData := amf0.Read(packet.Data)
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
							return err
						}
						i++
						fmt.Println("Metadata has been set")
						packet.CTX.Channel.Metadata = arr
						packet.CTX.IsMetadataSet = true
					}
				}
				i++
			}
		case 20:
			err = handleAMF0cmd(packet)
			return err
		}
	}
	return nil
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

	for {
		header, data, err := ctx.ReadPacket()
		if err != nil {
			log.Println("Error while reading Packet", err)
			return
		}

		err = handlePacket(ReceivedPacket{
			&ctx,
			&header,
			data,
		})
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
	}

}
