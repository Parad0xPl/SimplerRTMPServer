package amf0

// Read whole array and return all data encoded
func Read(input []byte, inputLength int) (error, []interface{}) {
	i := 0
	var parsedData []interface{}
	for i < inputLength {
		err, tmp, n := ReadAny(input[i:], inputLength-i)
		parsedData = append(parsedData, tmp)
		if err != nil {
			return err, parsedData
		}
		i += n
	}
	return nil, parsedData
}
