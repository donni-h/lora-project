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
	TypeRREQ MessageType = iota + '0'
	TypeRREP
	TypeRRER
	TypeData
)

func Unmarshal(data []byte) (Message, error) {
	var message Message
	if len(data) < 1 {
		return nil, fmt.Errorf("data can't be null")
	}
	switch MessageType(data[0]) {
	case TypeRREQ:
		message = &RREQ{}
	case TypeRREP:
		message = &RREP{}
	case TypeRRER:
		message = &RRER{}
	case TypeData:
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

/*
CompareSeqnums compares two sequence numbers.
It takes current and incoming sequence numbers as arguments, both of type int16.
The function returns true if the incoming sequence number is fresher (greater)
than or equal to the current sequence number, accounting for possible sequence number rollover.
Otherwise, it returns false.
*/
func CompareSeqnums(current int16, incoming int16) bool {
	return incoming-current >= 0
}
