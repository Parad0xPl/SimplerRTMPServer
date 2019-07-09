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
	conn net.Conn

	ChunkSize int
	initTime  uint32

	lastHeaderReceived Header
	lastHeaderSended   *Header

	SizeRead  int
	SizeWrote int

	Properties *map[string]interface{}

	ServerWindowAcknowledgement int
	ClientWindowAcknowledgement int
	PeerBandwidth               int
	PeerBandwidthType           int

	StreamID uint
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ctx ConnContext) read() (Header, []byte, error) {
	header, err := getHeader(&ctx)
	if err != nil {
		return Header{}, []byte{}, err
	}
	dataToRead := header.MessageLength
	offset := 0
	chunkLen := min(dataToRead, ctx.ChunkSize)
	body := make([]byte, header.MessageLength)
	tmp := make([]byte, chunkLen)
	for {
		_, err = ctx.conn.Read(tmp[:chunkLen])
		if err != nil {
			return Header{}, []byte{}, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = getHeader(&ctx)
		if err != nil {
			return Header{}, []byte{}, err
		}
	}
	return header, body, nil
}

// ReceivedPacket wrapper for all packet related data
type ReceivedPacket struct {
	ctx    *ConnContext
	header *Header
	data   []byte // Should work as long any handler won't try append to slice
}

func initCTX(conn net.Conn) ConnContext {
	return ConnContext{
		conn:                        conn,
		ChunkSize:                   128,
		initTime:                    getTime(),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Connection started: %s\n", conn.RemoteAddr().String())
	ctx := initCTX(conn)

	// Handle handshake
	err := handshake(conn)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// header, err := getHeader(&ctx)
		// ctx.SizeRead += header.Size
		// log.Println("Headers", header)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// data := make([]byte, header.MessageLength)
		// n, err := conn.Read(data)
		// ctx.SizeRead += n
		// if err != nil {
		// 	log.Println("Error while reading body", err)
		// 	return
		// }

		header, data, err := ctx.read()
		if err != nil {
			log.Println("Error while reading Packet", err)
			return
		}

		// //Magic byte fix
		// if n > 0x80 {
		// 	if data[0x80] == 0xc3 {
		// 		b := make([]byte, 1)
		// 		_, err := conn.Read(b)
		// 		if err != nil {
		// 			log.Println("Error while reading missing byte", err)
		// 			return
		// 		}
		// 		data = utils.Concat(data[:0x80], data[0x81:], b)
		// 	}
		// }

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
				return
			}
		} else {
			switch header.TypeID {
			case 4:
				// TODO handleUCM
				err = handleUCM(packet)
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
				handleAMF0cmd(packet)
				// TODO handle AMF0 command

			}
		}
	}

}
