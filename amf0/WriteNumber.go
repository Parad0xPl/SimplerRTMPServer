package amf0

import (
	"encoding/binary"
	"math"
)

// WriteNumber AMF0 format
func WriteNumber(n float64) []byte {
	b := make([]byte, 9)
	b[0] = 0
	binary.BigEndian.PutUint64(b[1:], uint64(math.Float64bits(n)))
	return b
}
