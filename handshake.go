package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

func handshake(c net.Conn) error {

	// Read c0
	c0 := make([]byte, 1)
	n, err := c.Read(c0)
	if err != nil {
		log.Println(err)
		return err
	}
	if n != 1 {
		return errors.New("Wrong length of chunk(c0): " + fmt.Sprint(n))
	}
	if c0[0] != 3 {
		return errors.New("Unsupported version: " + fmt.Sprint(c0[0]))
	}

	// Write s0
	_, err = c.Write([]byte{3})
	if err != nil {
		return err
	}

	// Make and send s1
	s1 := make([]byte, 1536)
	rand.Read(s1)
	binary.BigEndian.PutUint32(s1[0:4], 0)
	for i := 0; i < 4; i++ {
		s1[4+i] = 0
	}
	_, err = c.Write(s1)
	if err != nil {
		return err
	}

	// Read c1
	c1 := make([]byte, 1536)
	n, err = c.Read(c1)
	if err != nil {
		return err
	}
	if n != 1536 {
		return errors.New("Wrong length of handshake(c1): " + fmt.Sprint(n))
	}
	if !bytes.Equal(c1[4:8], []byte{0, 0, 0, 0}) {
		return errors.New("Zero section is not filled with zeros: " + fmt.Sprint(c1[4:8]))
	}

	// Send s2
	_, err = c.Write(c1)
	if err != nil {
		return err
	}

	// Read c2
	c2 := make([]byte, 1536)
	n, err = c.Read(c2)
	if err != nil {
		return err
	}
	if n != 1536 {
		return errors.New("Wrong length of handshake(c2): " + fmt.Sprint(n))
	}
	if !bytes.Equal(c2[4:8], []byte{0, 0, 0, 0}) {
		return errors.New("Zero section is not filled with zeros: " + fmt.Sprint(c2[4:8]))
	}

	// Check integrity of c2
	if !bytes.Equal(c2[0:4], s1[0:4]) {
		return errors.New("c2 time don't match s1`s\n" + fmt.Sprint(c2[0:4], s1[0:4]))
	}
	if !bytes.Equal(c2[4:8], c1[4:8]) {
		return errors.New("c2 time don't match s1`s\n" + fmt.Sprint(c2[4:8], c1[4:8]))
	}
	if !bytes.Equal(c2[8:], s1[8:]) {
		return errors.New("c2 random data don't match s1`s")
	}
	return nil
}
