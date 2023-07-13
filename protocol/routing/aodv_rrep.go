package routing

import (
	"log"
	"lora-project/protocol/messages"
)

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
