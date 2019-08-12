package amf0

// WriteAny in format
func WriteAny(raw interface{}) []byte {
	var buff []byte
	switch v := raw.(type) {
	case bool:
		buff = WriteBoolean(v)
	case string:
		buff = WriteString(v)
	case int:
		buff = WriteNumber(float64(v))
	case uint:
		buff = WriteNumber(float64(v))
	case float64:
		buff = WriteNumber(v)
	case nil:
		buff = []byte{5}
	case Undefined:
		buff = []byte{6}
	case ECMAArray:
		buff = WriteECAMArray(v.data)
	case map[string]interface{}:
		tmp := v
		if len(tmp) > 0 {
			buff = WriteObject(tmp)
		} else {
			buff = []byte{}
		}
	}
	return buff
}
