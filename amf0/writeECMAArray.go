package amf0

import (
	"SimpleRTMPServer/utils"
)

// ECMAArray _
type ECMAArray struct {
	data map[string]interface{}
}

// CreateECMAArray _
func CreateECMAArray(data map[string]interface{}) ECMAArray {
	return ECMAArray{
		data,
	}
}

// WriteECAMArray _
func WriteECAMArray(data map[string]interface{}) []byte {
	rawlen := len(data)
	parsed := make([][]byte, rawlen*2+3)
	parsed[0] = []byte{8}
	parsed[1] = utils.WriteInt(rawlen, 4)
	i := 2
	for key, val := range data {
		parsed[i] = WriteString(key)[1:]
		i++
		parsed[i] = WriteAny(val)
		i++
	}
	parsed[i] = []byte{0, 0, 9}
	return utils.Concat(parsed...)
}
