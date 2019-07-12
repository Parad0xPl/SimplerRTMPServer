package main

import (
	"SimpleRTMPServer/amf0"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func getTime() uint32 {
	return uint32(time.Now().UnixNano() / 1000)
}

type rawDataHandler func([]byte)

// ConnContext Structure for stream data and settings
type ConnContext struct {
	conn    net.Conn
	channel *ChannelObject

	readDumpFile  *os.File
	writeDumpFile *os.File

	audioHandler rawDataHandler
	videoHandler rawDataHandler

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

// Clear context before destructing
func (ctx *ConnContext) Clear() {
	if ctx.readDumpFile != nil {
		ctx.readDumpFile.Close()
	}
	if ctx.writeDumpFile != nil {
		ctx.writeDumpFile.Close()
	}
}

// Read Proxy for ctx.conn.Read
func (ctx ConnContext) Read(b []byte) (int, error) {
	n, err := ctx.conn.Read(b)
	if ctx.readDumpFile != nil {
		ctx.readDumpFile.Write(b[:n])
	}
	return n, err
}

// Write Proxy for ctx.conn.Write
func (ctx ConnContext) Write(b []byte) (int, error) {
	n, err := ctx.conn.Write(b)
	if ctx.writeDumpFile != nil {
		ctx.writeDumpFile.Write(b[:n])
	}
	return n, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ReadPacket _
func (ctx *ConnContext) ReadPacket() (Header, []byte, error) {
	header, err := getHeader(ctx)
	if err != nil {
		return Header{}, []byte{}, err
	}
	dataToRead := header.MessageLength
	offset := 0
	chunkLen := min(dataToRead, ctx.ChunkSize)
	body := make([]byte, header.MessageLength)
	tmp := make([]byte, chunkLen)
	for {
		_, err = ctx.Read(tmp[:chunkLen])
		if err != nil {
			return Header{}, []byte{}, err
		}
		offset += copy(body[offset:], tmp)
		chunkLen = min(dataToRead-offset, ctx.ChunkSize)
		if dataToRead-offset <= 0 {
			break
		}
		header, err = getHeader(ctx)
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

	ctx := ConnContext{
		conn:                        conn,
		ChunkSize:                   128,
		initTime:                    getTime(),
		ServerWindowAcknowledgement: 2500000,
		PeerBandwidth:               128,
	}

	if options.dumpfilein != "" &&
		options.dumpfileout != "" {
		n := options.dumpfilecounter
		options.dumpfilecounter++

		readfilename := fmt.Sprintf("%s.%d", options.dumpfilein, n)
		writefilename := fmt.Sprintf("%s.%d", options.dumpfileout, n)

		fmt.Printf(
			"Opening dump files\nInput data: %s\nOutput data: %s\n",
			readfilename,
			writefilename,
		)

		readfile, err := os.OpenFile(readfilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Read dump file")
		} else {
			ctx.readDumpFile = readfile
		}

		writefile, err := os.OpenFile(writefilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Couldn't open Write dump file")
		} else {
			ctx.writeDumpFile = writefile
		}
	}

	return ctx
}

func handler(conn net.Conn) {
	defer conn.Close()
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
