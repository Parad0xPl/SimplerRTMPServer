package utils

// ParseInt for protocol specification
func ParseInt(n, len int) []byte {
	v := make([]byte, len)
	for i := 0; i < len; i++ {
		v[i] = byte(n & 255)
		n = n >> 8
	}
	return v
}
