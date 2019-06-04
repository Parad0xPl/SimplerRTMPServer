package utils

// ReadInt for protocol specifitcation
func ReadInt(b []byte) int {
	tmp := 0
	for i := 0; i < len(b); i++ {
		tmp = tmp<<8 | int(b[i])
	}
	return tmp
}
