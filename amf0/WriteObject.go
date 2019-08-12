package amf0

import (
	"SimpleRTMPServer/utils"
	"bytes"
)

// WriteObject AMF0 format
func WriteObject(m map[string]interface{}) []byte {
	buff := new(bytes.Buffer)
	buff.Write([]byte{3})
	var tmp []byte
	for i, v := range m {
		tmp = []byte(i)
		buff.Write(utils.WriteIntBE(len(tmp), 2))
		buff.Write(tmp)
		buff.Write(WriteAny(v))
	}
	buff.Write([]byte{0, 0, 9})
	return buff.Bytes()
}
