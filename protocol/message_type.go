package protocol

import "fmt"

type MessageType uint8

const (
	MessageTypeRREQ MessageType = 0b00
	MessageTypeRREP MessageType = 0b01
	MessageTypeRRER MessageType = 0b10
	MessageTypeData MessageType = 0b11
)

func (m *MessageType) GetMessageType(data byte) Header {
	var header Header
	messageType := MessageType((data >> 6) & 0x03)
	switch messageType {
	case MessageTypeRREQ:
		header = &RREQHeader{}
	case MessageTypeRREP:
		header = &RREPHeader{}
	case MessageTypeRRER:
		header = &RRERHeader{}
	default:
		fmt.Println("unknown message type...")

	}
	return header

}
