package main

import (
	"flag"
)

// Options as flags` value
type Options struct {
	Port int // Ports used for RTMP

	FileConnectionInput  string // Attach file as net.Conn
	FileConnectionOutput string // Attach file as net.Conn

	DumpFileCounter   uint
	DumpInFnTemplate  string // Attach file as net.Conn
	DumpOutFnTemplate string // Attach file as net.Conn

	Analyze string // File with traffic to analise
}

func (opts *Options) init() {

	flag.IntVar(&opts.Port, "port", 1935, "RTMP port")

	flag.StringVar(&opts.FileConnectionInput, "cfi", "", "File attached as conn in (DEBUG function)")
	flag.StringVar(&opts.FileConnectionOutput, "cfo", "", "File attached as conn output (DEBUG function)")

	opts.DumpFileCounter = 0
	flag.StringVar(&opts.DumpInFnTemplate, "dfi", "", "File to dump connection (DEBUG function)")
	flag.StringVar(&opts.DumpOutFnTemplate, "dfo", "", "File to dump connection (DEBUG function)")

	flag.StringVar(&opts.Analyze, "analyze", "", "File to analyze (DEBUG function)")

	flag.Parse()
}
