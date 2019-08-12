package utils

// ToHex value between 0-15
func ToHex(b uint8) byte {
	if b >= 10 {
		return 'a' + b - 10
	}
	return '0' + b
}

// ByteToHex produce string of 8bit integer
func ByteToHex(b uint8) string {
	return string([]byte{ToHex(b >> 4), ToHex(b & 15)})
}
