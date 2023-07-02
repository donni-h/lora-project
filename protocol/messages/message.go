package messages

import "fmt"

type Message interface {
	Type() MessageType
	HeaderSize() int
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

type MessageType uint8

const (
	Type_RREQ MessageType = iota + '0'
	Type_RREP
	Type_RRER
	Type_Data
)

func Unmarshal(data []byte) (Message, error) {
	var message Message
	if len(data) < 1 {
		return nil, fmt.Errorf("data can't be null")
	}
	switch MessageType(data[0]) {
	case Type_RREQ:
		message = &RREQ{}
	case Type_RREP:
		message = &RREP{}
	case Type_RRER:
		message = &RRER{}
	case Type_Data:
		message = &Data{}
	default:
		return nil, fmt.Errorf("unknown message type: %s", string(data[0]))
	}

	err := message.Unmarshal(data[1:])
	if err != nil {
		return nil, err
	}

	return message, nil

}
