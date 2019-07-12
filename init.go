package main

import (
	"flag"
	"log"
)

// Options as flags` value
type Options struct {
	port int // Ports used for RTMP

	connfilein  string // Attach file as net.Conn
	connfileout string // Attach file as net.Conn

	dumpfilecounter uint
	dumpfilein      string // Attach file as net.Conn
	dumpfileout     string // Attach file as net.Conn
}

func initFlags() Options {
	log.Println("Initializing flags")
	options := Options{}

	flag.IntVar(&options.port, "port", 1935, "RTMP port")

	flag.StringVar(&options.connfilein, "connfilein", "", "File attached as conn in (DEBUG function)")
	flag.StringVar(&options.connfileout, "connfileout", "", "File attached as conn output (DEBUG function)")

	options.dumpfilecounter = 0
	flag.StringVar(&options.dumpfilein, "dumpfilein", "", "File to dump connection (DEBUG function)")
	flag.StringVar(&options.dumpfileout, "dumpfileout", "", "File to dump connection (DEBUG function)")

	flag.Parse()
	return options
}
