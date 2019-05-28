package main

import (
	"log"
	"net"
	"strconv"
)

var listner net.Listener

func serve(port int, handler func(net.Conn)) (err error) {
	// Listen on given port
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	listner = l
	defer l.Close()
	for {
		// Accept every connection
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		// pass handler to goroutine
		go handler(conn)
	}
}
