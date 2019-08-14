package amf0

import "errors"

// ReadBoolean data from byte array
func ReadBoolean(input []byte, inputLength int) (error, bool, int) {
	if inputLength < 1 {
		return errors.New("amf0 bool: need 1 byte"), false, 0
	}
	return nil, input[0] > 0, 1
}
