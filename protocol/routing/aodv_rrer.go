package routing

import "lora-project/protocol/messages"

func (a *AODV) handleRRER(rrer *messages.RRER, precursor messages.Address) {
	rt := a.routingTable
	existingEntry, exists := rt.GetEntry(rrer.UnreachDestinationAddress)

	// Discard message if it is not present in the current routing table
	if !exists || !messages.CompareSeqnums(existingEntry.SequenceNumber, rrer.UnreachDestinationSequence) ||
		existingEntry.NextHop != precursor {
		return
	}

	rt.DeleteEntry(rrer.UnreachDestinationAddress)

	a.sendToNextHop(rrer, a.broadcastAddress)
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
