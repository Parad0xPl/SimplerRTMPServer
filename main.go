package main

import (
	"SimpleRTMPServer/utils"
	"fmt"
	"log"
)

var options Options

func main() {
	fmt.Println("Starting SimpleRTMP Server")

	options = initFlags()

	initStrMan()

	if options.connfilein != "" &&
		options.connfileout != "" {
		faceconn, err := utils.OpenFileConn(options.connfilein,
			options.connfileout)
		if err != nil {
			log.Panicln(err)
		}
		handler(faceconn)

		return
	}

	err := serve(options.port, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
