package serial_handlers

import (
	"go.bug.st/serial"
	"lora-project/protocol/messages"
)

type UARTHandler struct {
	device            *serial.Port
	readCh            chan []byte
	writeCh           chan []byte
	errorCh           chan error
	stopCh            chan struct{}
	incomingMessageCh chan messages.Message
}

func NewUARTHandler(device *serial.Port) *UARTHandler {
	return &UARTHandler{
		device:            device,
		readCh:            make(chan []byte),
		writeCh:           make(chan []byte),
		errorCh:           make(chan error),
		stopCh:            make(chan struct{}),
		incomingMessageCh: make(chan messages.Message),
	}
}

func (handler *UARTHandler) readLoop() {

}

func (handler *UARTHandler) write() {

}
