package amf0

// Read whole array and return all data encoded
func Read(data []byte) []interface{} {
	i := 0
	dataLen := len(data)
	var parsedData []interface{}
	for i < dataLen {
		tmp, n := ReadAny(data[i:])
		parsedData = append(parsedData, tmp)
		i += n
	}
	return parsedData
}
