package utils

// WriteInt for protocol specification
func WriteInt(n, len int) []byte {
	v := make([]byte, len)
	if n >= (1 << (uint(len) * 8)) {
		n = (1 << (uint(len) * 8)) - 1
	}
	for i := 0; i < len; i++ {
		v[i] = byte(n & 255)
		n = n >> 8
	}
	return v
}
