package utils

// ReadUint for protocol specification
func ReadUint(b []byte) uint {
	tmp := uint(0)
	blen := len(b)
	for i := 0; i < blen; i++ {
		tmp = tmp<<8 | uint(b[i])
	}
	return tmp
}
