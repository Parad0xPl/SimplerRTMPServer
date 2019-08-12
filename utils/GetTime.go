package utils

import "time"

// GetTime in milliseconds
func GetTime() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}
