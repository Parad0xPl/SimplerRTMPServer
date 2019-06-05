package amf3

import (
	"SimpleRTMPServer/amf0"
)

// ReadNumber amf3 encoded data
func ReadNumber(data []byte) (float64, int) {
	return amf0.ReadNumber(data)
}
