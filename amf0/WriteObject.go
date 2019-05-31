package amf0

import "bytes"

// WriteObject AMF0 format
func WriteObject(m map[string]interface{}) []byte {
	buff := new(bytes.Buffer)
	buff.Write([]byte{3})
	for i, v := range m {
		buff.Write([]byte(i))
		buff.Write(WriteAny(v))
	}
	buff.Write([]byte{0, 0, 9})
	return buff.Bytes()
}
