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
		buff = WriteObject(v.(map[string]interface{}))
	}
	return buff
}
