package messages

import "fmt"

type Message interface {
	HeaderSize() int
	Marshal() ([]byte, error)
	unmarshal(data []byte) error
}

func UnmarshalHeader(data []byte) (Message, error) {
	var message Message

	switch data[0] {
	case '0':
		message = &RREQ{}
	case '1':
		message = &RREP{}
	case '2':
		message = &RRER{}
	case '3':
		message = &Data{}
	case '4':
		message = &DataAck{}
	default:
		return nil, fmt.Errorf("unknown message type: %s", string(data[0]))
	}

	err := message.unmarshal(data)
	if err != nil {
		return nil, err
	}

	return message, nil

}
