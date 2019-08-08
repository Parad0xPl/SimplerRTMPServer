package main

import "time"

func getTime() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}
