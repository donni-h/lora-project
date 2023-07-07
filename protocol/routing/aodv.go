package routing

import (
	"log"
	"lora-project/protocol/messages"
	"lora-project/serial_handlers"
)

type AODV struct {
	handler          *serial_handlers.ATHandler
	routingTable     Table
	idTable          RREQIDTable
	rreqID           uint16
	seqNum           int16
	currentAddress   messages.Address
	messageQueue     chan messages.Message
	broadcastAddress messages.Address
}

func NewAODV(atHandler *serial_handlers.ATHandler) *AODV {
	var addr messages.Address
	err := addr.UnmarshalText([]byte("4761"))
	if err != nil {
		log.Fatal(err)
	}
	var broadcast messages.Address
	err = broadcast.UnmarshalText([]byte("FFFF"))
	if err != nil {

	}
	return &AODV{
		handler:          atHandler,
		routingTable:     NewTable(),
		idTable:          NewRREQIDTable(),
		rreqID:           0,
		seqNum:           0,
		currentAddress:   addr,
		messageQueue:     make(chan messages.Message, 10),
		broadcastAddress: broadcast,
	}
}
func (a *AODV) Run() {
	go a.processIncomingMessages()
	go a.helloRoutine()
}

func (a *AODV) processIncomingMessages() {
	for {
		select {
		case msg := <-a.handler.MessageChan:
			a.processMessage(msg)
		case err := <-a.handler.ErrorChan:
			a.processError(err)
		}
	}
}
