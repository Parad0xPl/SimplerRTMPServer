package amf0

import (
	"SimpleRTMPServer/utils"
)

// WriteString AMF0 format
func WriteString(str string) []byte {
	tmp := []byte(str)
	return utils.Concat([]byte{2}, utils.WriteIntBE(len(tmp), 2), tmp)
}
