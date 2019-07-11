package amf0

import (
	"SimpleRTMPServer/utils"
)

// ECMAArray _
type ECMAArray struct {
	data map[string]interface{}
}

// WriteECAMArray _
func WriteECAMArray(data map[string]interface{}) []byte {
	rawlen := len(data)
	parsed := make([][]byte, rawlen*2+2)
	parsed[0] = utils.WriteInt(rawlen, 4)
	i := 1
	for key, val := range data {
		parsed[i] = WriteString(key)
		i++
		parsed[i] = WriteAny(val)
		i++
	}
	parsed[i] = []byte{0, 0, 9}
	return utils.Concat(parsed...)
}
