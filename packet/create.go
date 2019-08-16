package packet

import (
	"SimpleRTMPServer/amf0"
)

type create struct{}

// Create a packet
var Create create

// Prototype packet skeleton
type Prototype struct {
	Head *Header
	Body []byte
}

func (create) CommandMessage(raw []interface{}) Prototype {
	head := &Header{
		MessageTypeID: 20,
		ChunkStreamID: 3,
	}
	body := amf0.Write(raw)
	return Prototype{head, body}
}

type CommandArgs struct {
	TransactionID int
	Properties    interface{}
	Information   interface{}
}

func (create) ResultMessage(args CommandArgs) Prototype {
	raw := make([]interface{}, 4)
	raw[0] = "_result"
	raw[1] = args.TransactionID
	raw[2] = args.Properties
	raw[3] = args.Information
	packet := Create.CommandMessage(raw)
	return packet
}

func (create) ErrorMessage(args CommandArgs) Prototype {
	raw := make([]interface{}, 4)
	raw[0] = "_error"
	raw[1] = args.TransactionID
	raw[2] = args.Properties
	raw[3] = args.Information
	packet := Create.CommandMessage(raw)
	return packet
}

func (create) OnStatusMessage(level, code, desc string) Prototype {
	raw := make([]interface{}, 4)
	raw[0] = "onStatus"
	raw[1] = 0
	raw[2] = nil
	raw[3] = map[string]interface{}{
		"level":       level,
		"code":        code,
		"description": desc,
	}
	packet := Create.CommandMessage(raw)
	return packet
}

func (create) AudioData(data []byte) Prototype {
	head := &Header{
		MessageTypeID: 8,
	}
	return Prototype{head, data}
}

func (create) VideoData(data []byte) Prototype {
	head := &Header{
		MessageTypeID: 9,
	}
	return Prototype{head, data}
}

func (create) AMF0Data(data []interface{}) Prototype {
	head := &Header{
		MessageTypeID: 0x12,
	}
	body := amf0.Write(data)
	return Prototype{head, body}
}
