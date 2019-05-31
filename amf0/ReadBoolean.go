package amf0

// ReadBoolean data from byte array
func ReadBoolean(data []byte) (bool, int) {
	return data[0] > 0, 1
}
