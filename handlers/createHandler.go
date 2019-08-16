package handlers

import (
	"SimpleRTMPServer/connCTX"
	"SimpleRTMPServer/packet"
	"log"
)

func createAudioHandler(ctx *connCTX.ConnContext) connCTX.RawDataHandler {
	return func(data []byte) {
		pkt := packet.Create.AudioData(data)
		for _, v := range ctx.Channel.Subscribed {
			v.SendPacket(pkt)
			log.Println("Audio data send to client")
		}
	}
}

func createVideoHandler(ctx *connCTX.ConnContext) connCTX.RawDataHandler {
	return func(data []byte) {
		pkt := packet.Create.VideoData(data)
		for _, v := range ctx.Channel.Subscribed {
			v.SendPacket(pkt)
			log.Println("Video data send to client")
		}
	}
}
