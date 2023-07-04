package protocol

import (
	"log"
	"lora-project/protocol/messages"
	"lora-project/protocol/routing"
	"lora-project/serial_handlers"
)

type AODV struct {
	handler        *serial_handlers.ATHandler
	routingTable   routing.Table
	idTable        routing.RREQIDTable
	rreqID         uint16
	currentAddress messages.Address
}

func NewAODV(atHandler *serial_handlers.ATHandler) *AODV {
	var addr messages.Address
	err := addr.UnmarshalText([]byte("4761"))
	if err != nil {
		log.Fatal(err)
	}
	return &AODV{
		handler:        atHandler,
		routingTable:   routing.NewTable(),
		idTable:        routing.NewRREQIDTable(),
		rreqID:         0,
		currentAddress: addr,
	}
}
