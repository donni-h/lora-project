package routing

import (
	"fmt"
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
			a.handleData(data, pre)
		}

	default:
		log.Printf("Received unknown message type: %v\n", msg.Type())
	}
}

func (a *AODV) handleRREQ(rreq *messages.RREQ, precursor messages.Address) {
	fmt.Println("Received RREQ:", rreq)
	fmt.Println("Precursor:", precursor.String())
}

func (a *AODV) handleRREP(rrep *messages.RREP, precursor messages.Address) {
	rt := a.routingTable

	existingEntry, exists := rt.GetEntry(rrep.DestinationAddress)

	if !exists || messages.CompareSeqnums(existingEntry.SequenceNumber, rrep.DestinationSequenceNum) {
		rt.AddOrUpdateEntry(
			rrep.DestinationAddress,
			precursor,
			rrep.HopCount+1,
			[]messages.Address{},
			rrep.DestinationSequenceNum,
		)
	}

	if rrep.OriginatorAddress == a.currentAddress {
		return
	}

	originatorEntry, exists := rt.GetEntry(rrep.OriginatorAddress)

	if !exists {
		// If there is no known route to the originator, we just drop the packet
		return
	}

	nextHop := originatorEntry.NextHop
	rrep.HopCount += 1
	a.sendToNextHop(rrep, nextHop)
}

func (a *AODV) handleRRER(rrer *messages.RRER) {
	rt := a.routingTable
	existingEntry, exists := rt.GetEntry(rrer.UnreachDestinationAddress)

	// Discard message if it is not present in the current routing table
	if !exists || !messages.CompareSeqnums(existingEntry.SequenceNumber, rrer.UnreachDestinationSequence) {
		return
	}

	rt.DeleteEntry(rrer.UnreachDestinationAddress)

	a.sendToNextHop(rrer, a.broadcastAddress)
}

func (a *AODV) handleData(data *messages.Data, precursor messages.Address) {
	fmt.Println("Received Data:", data)
	fmt.Println("Precursor:", precursor.String())
}

func (a *AODV) sendToNextHop(msg messages.Message, nextHop messages.Address) {
	a.handler.SetTargetAddress(nextHop)
	a.handler.SendMessage(msg)
}
