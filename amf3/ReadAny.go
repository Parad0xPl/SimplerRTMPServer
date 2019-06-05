package amf3

import (
	"log"
)

// ReadAny amf3 encoded data
func ReadAny(data []byte) (interface{}, int) {
	var tmp interface{}
	var n int
	switch data[0] {
	case 0:
		return nil, 1
	case 1:
		return nil, 1
	case 2:
		return false, 1
	case 3:
		return true, 1
	case 4:
		tmp, n = ReadInt(data[1:])
	case 5:
		tmp, n = ReadNumber(data[1:])
	// case 6:
	//String read
	// case 0x0a:
	//Object read
	default:
		log.Println("Unsupported type in amf3 decoder")
		return nil, 1
	}
	return tmp, n + 1
}
