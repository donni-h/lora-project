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

func (a *AODV) handleRREQ(rreq *messages.RREQ, precursor messages.Address) {
	rt := a.routingTable
	if a.idTable.hasID(rreq.OriginatorAddress, rreq.RREQID) || rreq.OriginatorAddress == a.currentAddress {
		return
	}

	existingEntry, exists := rt.GetEntry(rreq.OriginatorAddress)

	if !exists || messages.CompareSeqnums(existingEntry.SequenceNumber, rreq.OriginatorSequenceNum) {
		rt.AddOrUpdateEntry(
			rreq.OriginatorAddress,
			precursor,
			rreq.HopCount+1,
			[]messages.Address{},
			rreq.OriginatorSequenceNum,
		)
	}
	a.idTable.AddID(rreq.OriginatorAddress, rreq.RREQID)

	if rreq.DestinationAddress == a.currentAddress {
		a.generateRREP(rreq.OriginatorAddress, rreq.DestinationAddress)
		return
	}

	_, exists = rt.GetEntry(rreq.DestinationAddress)

	if exists {
		a.generateRREP(rreq.OriginatorAddress, rreq.DestinationAddress)
	}

	rreq.HopCount += 1
	a.sendToNextHop(rreq, a.broadcastAddress)
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

	if rrep.OriginatorAddress == a.currentAddress || rrep.OriginatorAddress == a.broadcastAddress {
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

func (a *AODV) handleData(data *messages.Data) {
	if messages.CompareSeqnums(a.seqNum, data.DataSequenceNumber) {
		a.seqNum = data.DataSequenceNumber
	}
	if data.DestinationAddress == a.currentAddress {
		a.incomingDataQueue <- data
		return
	}

	route, exists := a.routingTable.GetEntry(data.DestinationAddress)
	if !exists {
		a.generateRRER(data.DestinationAddress)
		return
	}
	a.sendToNextHop(data, route.NextHop)
}

func (a *AODV) sendToNextHop(msg messages.Message, nextHop messages.Address) {
	a.handler.SetTargetAddress(nextHop)
	a.handler.SendMessage(msg)
}

func (a *AODV) generateRREP(originator messages.Address, destination messages.Address) {
	rt := a.routingTable
	destEntry, exists := rt.GetEntry(destination)

	if !exists {
		log.Printf("No route to destination address: %s\n", destination.String())
		return

	}

	rrep := &messages.RREP{
		HopCount:               destEntry.HopCount,
		DestinationAddress:     destination,
		DestinationSequenceNum: destEntry.SequenceNumber,
		OriginatorAddress:      originator,
	}
	originatorEntry, exists := rt.GetEntry(originator)

	if !exists {
		log.Printf("No route to originator address: %s\n", originator.String())
		a.generateRRER(originator)
	}

	a.sendToNextHop(rrep, originatorEntry.NextHop)
}

func (a *AODV) sendData(payload string, destination messages.Address) {
	a.seqNum++
	data := &messages.Data{
		DestinationAddress: destination,
		OriginatorAddress:  a.currentAddress,
		DataSequenceNumber: a.seqNum,
		Payload:            []byte(payload),
	}

	entry, exists := a.routingTable.GetEntry(destination)
	if !exists {
		_, exists = a.dataQueue.conds[destination]
		a.dataQueue.Push(data)

		if !exists {
			a.generateRREQ(destination)
			go a.dataQueue.Pop(destination)
		}
		return
	}
	a.sendToNextHop(data, entry.NextHop)
}

func (a *AODV) generateRREQ(address messages.Address) {
	a.rreqID += 1
	a.seqNum += 1

	rreq := &messages.RREQ{
		UFlag:                  true,
		HopCount:               0,
		RREQID:                 a.rreqID,
		DestinationAddress:     address,
		DestinationSequenceNum: 0,
		OriginatorAddress:      a.currentAddress,
		OriginatorSequenceNum:  a.seqNum,
	}

	a.sendToNextHop(rreq, a.broadcastAddress)
}

func (a *AODV) queueData(destination messages.Address) {
	go func() {
		pendingData := a.dataQueue.Pop(destination)

		for _, msg := range pendingData {
			entry, ok := a.routingTable.GetEntry(destination)

			if !ok {
				log.Println("Couldn't send message, no route to:" + destination.String())
				return
			}

			a.sendToNextHop(msg, entry.NextHop)
		}
	}()
}

func (a *AODV) generateRRER(address messages.Address) {
	entry, exists := a.routingTable.GetEntry(address)

	if !exists {
		return
	}
	rrer := &messages.RRER{
		UnreachDestinationAddress:  address,
		UnreachDestinationSequence: entry.SequenceNumber,
	}
	a.handler.SendMessage(rrer)
}
