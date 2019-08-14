package amf0

import (
	"encoding/binary"
	"errors"
	"math"
)

// ReadNumber read amf0 number from byte array
func ReadNumber(input []byte, inputLength int) (error, float64, int) {
	if inputLength < 8 {
		return errors.New("amf0 number: need 8 bytes"), 0, 0
	}
	bits := binary.BigEndian.Uint64(input)
	float := math.Float64frombits(bits)
	return nil, float, 8
}
