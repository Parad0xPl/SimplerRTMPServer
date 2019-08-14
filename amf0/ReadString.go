package amf0

import (
	"encoding/binary"
	"errors"
)

// ReadString read amf0 string from byte array
func ReadString(input []byte, inputLength int) (error, string, int) {
	if inputLength < 2 {
		return errors.New("amf0 string: no length"), "", 0
	}
	dataLen := int(binary.BigEndian.Uint16(input))
	if inputLength < 2+dataLen {
		return errors.New("amf0 string: inputLength is less then string length"), "", 0
	}
	return nil, string(input[2 : dataLen+2]), dataLen + 2
}
