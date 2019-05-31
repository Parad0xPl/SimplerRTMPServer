package amf0

import "bytes"

// Write message in AMF0 format
func Write(a []interface{}) []byte {
	buff := new(bytes.Buffer)
	for _, v := range a {
		buff.Write(WriteAny(v))
	}
	return buff.Bytes()
}
