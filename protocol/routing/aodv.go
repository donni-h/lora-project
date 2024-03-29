package routing

import (
	"log"
	"lora-project/protocol/messages"
	"lora-project/serial_handlers"
)

type AODV struct {
	handler           *serial_handlers.ATHandler
	routingTable      *Table
	idTable           *RREQIDTable
	rreqID            uint16
	seqNum            int16
	currentAddress    messages.Address
	IncomingDataQueue chan messages.Data
	broadcastAddress  messages.Address
	dataQueue         *DataQueue
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
		handler:           atHandler,
		routingTable:      NewTable(),
		idTable:           NewRREQIDTable(),
		rreqID:            0,
		seqNum:            0,
		currentAddress:    addr,
		IncomingDataQueue: make(chan messages.Data, 10),
		broadcastAddress:  broadcast,
		dataQueue:         NewDataQueue(),
	}
}
func (a *AODV) Run() {
	str := []string{"433920000", "5", "9", "7", "4", "1", "0", "0", "0", "0", "3000", "8", "8"}
	go a.processIncomingMessages()
	a.handler.SetOwnAddress(a.currentAddress)
	a.handler.SetTargetAddress(a.broadcastAddress)
	a.handler.Configure(str)
	a.handler.SetReceive()
	go a.helloRoutine()
	a.StartExpirationWorker()
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
