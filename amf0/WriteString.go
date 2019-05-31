package amf0

// WriteString AMF0 format
func WriteString(str string) []byte {
	return append([]byte{2}, []byte(str)...)
}
