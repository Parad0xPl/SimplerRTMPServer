package amf0

import (
	"SimpleRTMPServer/utils"
	"bytes"
	"log"
)

// ReadECMAArray _
func ReadECMAArray(data []byte) (map[string]interface{}, int) {
	offset := 0
	alen := utils.ReadInt(data[0:4])
	offset += 4
	i := 0
	ret := make(map[string]interface{})
	for i < alen {
		prop, n := ReadString(data[offset:])
		offset += n
		item, n := ReadAny(data[offset:])
		offset += n
		ret[prop] = item
		i++
	}
	if bytes.Compare(data[offset:offset+3], []byte{0x0, 0x0, 0x09}) != 0 {
		log.Println("ECMA Array don't end with EoB marker")
	}
	offset += 3
	return ret, offset
}
