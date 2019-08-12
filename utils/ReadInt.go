package utils

// ReadInt for protocol specification
func ReadInt(b []byte) int {
	tmp := 0
	bLen := len(b)
	for i := 0; i < bLen; i++ {
		tmp = tmp<<8 | int(b[i])
	}
	return tmp
}
