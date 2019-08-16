package main

import (
	"SimpleRTMPServer/amf0"
	"SimpleRTMPServer/handlers"
	"SimpleRTMPServer/packet"
	"fmt"
	"log"
	"net"
)

func handlePacket(recvPacket handlers.ReceivedPacket) error {
	var err error
	if recvPacket.Header.ChunkStreamID == 2 && recvPacket.Header.MessageStreamID == 0 {
		// Take effect when received
		err = handlers.PCM(recvPacket)
		return err
	} else {
		switch recvPacket.Header.MessageTypeID {
		case 8:
			if recvPacket.CTX.AudioHandler != nil {
				log.Println("Audio data received. Len:", len(recvPacket.Data))
				recvPacket.CTX.AudioHandler(recvPacket.Data)
			}
			return nil
		case 9:
			if recvPacket.CTX.VideoHandler != nil {
				log.Println("Video data received. Len:", len(recvPacket.Data))
				recvPacket.CTX.VideoHandler(recvPacket.Data)
			}
			return nil
		case 15:
			// TODO handle AMF3 data
		case 17:
			// TODO handle AMF3 command
		case 18:
			log.Println("AMF0 data received")
			err, parsedData := amf0.Read(recvPacket.Data, len(recvPacket.Data))
			if err != nil {
				return err
			}
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
							return nil
						}
						i++
						fmt.Println("Metadata has been set")
						for _, cn := range recvPacket.CTX.Channel.Subscribed {
							pkt := packet.Create.AMF0Data([]interface{}{
								"onMetaData",
								amf0.CreateECMAArray(arr),
							})
							cn.SendPacket(pkt)
						}
						recvPacket.CTX.Channel.Metadata = arr
						recvPacket.CTX.IsMetadataSet = true
					}
				}
				i++
			}
		case 20:
			err = handlers.AMF0cmd(recvPacket)
			return err
		}
	}
	return nil
}

func handler(conn net.Conn) {
	fmt.Printf("Connection started: %s\n", conn.RemoteAddr().String())
	ctx := NewCTX(conn)
	defer ctx.Clear()

	// Handle handshake
	err := handshake(ctx)
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

		err = handlePacket(handlers.ReceivedPacket{
			ctx,
			serverInstance,
			&header,
			data,
		})
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
	}

}
