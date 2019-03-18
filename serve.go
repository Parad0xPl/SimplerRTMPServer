package main

import (
	"log"
	"net"
	"strconv"
)

var listner net.Listener

func serve(port int, handler func(net.Conn)) (err error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	listner = l
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		go handler(conn)
	}
}
