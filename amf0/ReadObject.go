package amf0

import "bytes"

// ReadObject return map
func ReadObject(input []byte, inputLength int) (error, map[string]interface{}, int) {
	var parsedData map[string]interface{}
	var err error = nil
	parsedData = make(map[string]interface{})
	i := 0
	for {
		if bytes.Compare(input[i:], []byte{0, 0, 9}) == 0 {
			i += 3
			break
		}
		var key string
		var n, tmpLen int
		var tmp interface{}
		err, key, n = ReadString(input[i:], inputLength-i)
		if err != nil {
			return err, parsedData, 0
		}
		i += n
		err, tmp, tmpLen = ReadAny(input[i:], inputLength-i)
		i += tmpLen
		parsedData[key] = tmp
		if err != nil {
			return err, parsedData, 0
		}
	}
	return err, parsedData, i
}
