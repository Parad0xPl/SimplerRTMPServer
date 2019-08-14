package amf0

import (
	"SimpleRTMPServer/utils"
	"bytes"
	"log"
)

// ReadECMAArray _
func ReadECMAArray(input []byte, inputLength int) (error, map[string]interface{}, int) {
	var err error = nil
	offset := 0
	aLen := utils.ReadInt(input[0:4])
	offset += 4
	i := 0
	ret := make(map[string]interface{})
	for i < aLen {
		var prop string
		var n int
		var item interface{}
		err, prop, n = ReadString(input[offset:], inputLength-offset)
		offset += n
		if err != nil {
			return err, ret, offset
		}
		err, item, n = ReadAny(input[offset:], inputLength-offset)
		offset += n
		ret[prop] = item
		if err != nil {
			return err, ret, offset
		}
		i++
	}
	if bytes.Compare(input[offset:offset+3], []byte{0x0, 0x0, 0x09}) != 0 {
		log.Println("ECMA Array don't end with EoB marker")
	}
	offset += 3
	return err, ret, offset
}
