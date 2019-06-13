package amf0

// WriteAny in format
func WriteAny(v interface{}) []byte {
	var buff []byte
	switch v.(type) {
	case bool:
		buff = WriteBoolean(v.(bool))
	case string:
		buff = WriteString(v.(string))
	case float64:
		buff = WriteNumber(v.(float64))
	case nil:
		buff = []byte{5}
	case map[string]interface{}:
		tmp := v.(map[string]interface{})
		if len(tmp) > 0 {
			buff = WriteObject(tmp)
		} else {
			buff = []byte{}
		}
	}
	return buff
}
