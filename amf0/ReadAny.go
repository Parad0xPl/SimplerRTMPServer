package amf0

import "errors"

// ReadAny data from byte array
func ReadAny(input []byte, inputLength int) (error, interface{}, int) {
	var tmp interface{}
	var err error = nil
	var n int
	if inputLength < 1 {
		return errors.New("input too short"), nil, 0
	}
	switch input[0] {
	case 0x00:
		err, tmp, n = ReadNumber(input[1:], inputLength-1)
	case 0x01:
		err, tmp, n = ReadBoolean(input[1:], inputLength-1)
	case 0x02:
		err, tmp, n = ReadString(input[1:], inputLength-1)
	case 0x03:
		err, tmp, n = ReadObject(input[1:], inputLength-1)
	case 0x05:
		tmp = nil
		n = 0
	case 0x06:
		tmp = Undefined{}
		n = 0
	case 0x08:
		err, tmp, n = ReadECMAArray(input[1:], inputLength-1)
	default:
		err = errors.New("unsupported AMF0 type")
		return err, nil, 1
	}
	return err, tmp, n + 1
}
