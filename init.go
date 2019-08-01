package main

import (
	"flag"
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

func (opts *Options) init() {

	flag.IntVar(&opts.port, "port", 1935, "RTMP port")

	flag.StringVar(&opts.connfilein, "cfi", "", "File attached as conn in (DEBUG function)")
	flag.StringVar(&opts.connfileout, "cfo", "", "File attached as conn output (DEBUG function)")

	opts.dumpfilecounter = 0
	flag.StringVar(&opts.dumpfilein, "dfi", "", "File to dump connection (DEBUG function)")
	flag.StringVar(&opts.dumpfileout, "dfo", "", "File to dump connection (DEBUG function)")

	flag.Parse()
}
