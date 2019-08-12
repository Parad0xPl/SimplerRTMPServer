package utils

// CountBits return amount of bits in v
func CountBits(v uint64) uint {
	counter := uint(0)
	for v != 0 {
		v &= v - 1
		counter++
	}
	return counter
}

// RotateBitsRight rotate right a bits of v
func RotateBitsRight(v uint64, a uint8) uint64 {
	a &= 64 - 1
	return (v << (64 - a)) | (v >> a)
}

// RotateBitsLeft rotate left a bits of v
func RotateBitsLeft(v uint64, a uint8) uint64 {
	a &= 64 - 1
	return RotateBitsRight(v, 64-a)
}
