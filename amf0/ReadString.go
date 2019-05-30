package amf0

import "encoding/binary"

// ReadString read amf0 string from byte array
func ReadString(data []byte) (string, int) {
	len := int(binary.BigEndian.Uint16(data))
	return string(data[2 : len+2]), len + 2
}