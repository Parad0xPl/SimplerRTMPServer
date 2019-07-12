package main

func createAudioHandler(ctx *ConnContext) rawDataHandler {
	return func(data []byte) {
		pkt := Create.audiodata(int(ctx.StreamID), data)
		for _, v := range ctx.channel.subscribed {
			sendPacket(v, pkt)
		}
	}
}

func createVideoHandler(ctx *ConnContext) rawDataHandler {
	return func(data []byte) {
		pkt := Create.videodata(int(ctx.StreamID), data)
		for _, v := range ctx.channel.subscribed {
			sendPacket(v, pkt)
		}
	}
}
