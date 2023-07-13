package routing

import (
	"log"
	"lora-project/protocol/messages"
	"lora-project/serial_handlers"
)

func (a *AODV) processMessage(event serial_handlers.MessageEvent) {
	msg, pre := event.Message, event.Precursor
	switch msg.Type() {
	case messages.TypeRREQ:
		if rreq, ok := msg.(*messages.RREQ); ok {
			a.handleRREQ(rreq, pre)
		}
	case messages.TypeRREP:
		if rrep, ok := msg.(*messages.RREP); ok {
			a.handleRREP(rrep, pre)
		}
	case messages.TypeRRER:
		if rrer, ok := msg.(*messages.RRER); ok {
			a.handleRRER(rrer)
		}
	case messages.TypeData:
		if data, ok := msg.(*messages.Data); ok {
			a.handleData(data)
		}

	default:
		log.Printf("Received unknown message type: %v\n", msg.Type())
	}
}

func (a *AODV) sendToNextHop(msg messages.Message, nextHop messages.Address) {
	a.handler.SetTargetAddress(nextHop)
	a.handler.SendMessage(msg)
}
