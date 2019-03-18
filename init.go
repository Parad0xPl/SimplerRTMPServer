package main

import (
	"flag"
	"log"
)

// Options as flags` value
type Options struct {
	port int // Ports used for RTMP
}

func initFlags() Options {
	log.Println("Initializing flags")
	options := Options{}

	flag.IntVar(&options.port, "port", 1935, "RTMP port")

	flag.Parse()
	return options
}
