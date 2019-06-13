package utils

// ReadInt for protocol specifitcation
func ReadInt(b []byte) int {
	tmp := 0
	blen := len(b)
	for i := 0; i < blen; i++ {
		tmp = tmp<<8 | int(b[i])
	}
	return tmp
}
