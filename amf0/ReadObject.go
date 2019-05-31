package amf0

import "bytes"

// ReadObject return map
func ReadObject(data []byte) (map[string]interface{}, int) {
	var parsedData map[string]interface{}
	parsedData = make(map[string]interface{})
	i := 0
	for {
		if bytes.Compare(data[i:], []byte{0, 0, 9}) == 0 {
			i += 3
			break
		}
		key, n := ReadString(data[i:])
		i += n
		tmp, tmplen := ReadAny(data[i:])
		i += tmplen
		parsedData[key] = tmp
	}
	return parsedData, i
}
