package main

import (
	"SimpleRTMPServer/utils"
	"fmt"
	"log"
)

var options Options

func main() {
	fmt.Println("Starting SimpleRTMP Server")

	options.init()
	initStrMan()

	if options.FileConnectionInput != "" &&
		options.FileConnectionOutput != "" {
		faceconn, err := utils.OpenFileConn(options.FileConnectionInput,
			options.FileConnectionOutput)
		if err != nil {
			log.Panicln(err)
		}
		handler(&faceconn)

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
