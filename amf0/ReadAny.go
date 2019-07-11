package amf0

// ReadAny data from byte array
func ReadAny(data []byte) (interface{}, int) {
	var tmp interface{}
	var n int
	switch data[0] {
	case 0:
		tmp, n = ReadNumber(data[1:])
	case 1:
		tmp, n = ReadBoolean(data[1:])
	case 2:
		tmp, n = ReadString(data[1:])
	case 3:
		tmp, n = ReadObject(data[1:])
	case 5:
		tmp = nil
		n = 0
	case 8:
		tmp, n = ReadECMAArray(data[1:])
	default:
		return nil, 1
	}
	return tmp, n + 1
}
