package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting SimpleRTMP Server")

	options := initFlags()

	err := serve(options.port, handler)
	if err != nil {
		log.Fatalln(err)
	}
}