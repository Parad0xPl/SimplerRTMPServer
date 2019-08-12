package main

import (
	"SimpleRTMPServer/utils"
	"fmt"
	"log"
)

var options Options

func init() {
	options.init()
}

func main() {
	fmt.Println("Starting SimpleRTMP Server")

	initServerInstance()

	if options.FileConnectionInput != "" &&
		options.FileConnectionOutput != "" {
		faceConn, err := utils.OpenFileConn(options.FileConnectionInput,
			options.FileConnectionOutput)
		if err != nil {
			log.Panicln(err)
		}
		handler(&faceConn)

		return
	}

	if options.Analyze != "" {
		analyze(options.Analyze)
		return
	}

	err := serve(options.Port, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
