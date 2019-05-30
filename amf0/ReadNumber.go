package amf0

import (
	"encoding/binary"
	"math"
)

// ReadNumber read amf0 number from byte array
func ReadNumber(data []byte) (float64, int) {
	bits := binary.BigEndian.Uint64(data)
	float := math.Float64frombits(bits)
	return float, 8
}
