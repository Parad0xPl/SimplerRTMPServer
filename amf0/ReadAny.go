package amf0

// ReadAny data from byte array
func ReadAny(data []byte) (interface{}, int) {
	switch data[0] {
	case 0:
		tmp, n := ReadNumber(data[1:])
		return tmp, n + 1
	case 2:
		str, n := ReadString(data[1:])
		return str, n + 1

	default:
		return nil, 1
	}
}
