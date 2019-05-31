package amf0

// WriteBoolean AMF0 format
func WriteBoolean(b bool) []byte {
	if b {
		return []byte{1, 1}
	}
	return []byte{1, 0}
}
