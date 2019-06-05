package amf3

// ReadInt amf3 encoded data
func ReadInt(data []byte) (int, int) {
	tmp := 0
	i := 0
	for {
		tmp = int(data[i]) & 0x7f
		i++
		if tmp&0x8 > 0 {
			break
		}
		tmp = tmp << 7
	}
	return tmp, i
}
