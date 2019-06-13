package amf0

import (
	"SimpleRTMPServer/utils"
)

// WriteString AMF0 format
func WriteString(str string) []byte {
	tmp := []byte(str)
	begin := append([]byte{2}, utils.WriteInt(len(tmp), 2)...)
	return append(begin, tmp...)
}
