package main

import "log"

func createAudioHandler(ctx *ConnContext) RawDataHandler {
	return func(data []byte) {
		pkt := Create.AudioData(data)
		for _, v := range ctx.Channel.Subscribed {
			v.SendPacket(pkt)
			log.Println("Audio data send to client")
		}
	}
}

func createVideoHandler(ctx *ConnContext) RawDataHandler {
	return func(data []byte) {
		pkt := Create.VideoData(data)
		for _, v := range ctx.Channel.Subscribed {
			v.SendPacket(pkt)
			log.Println("Video data send to client")
		}
	}
}
