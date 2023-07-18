package routing

import (
	"fmt"
	"log"
	"lora-project/protocol/messages"
	"time"
)

func (a *AODV) generateRREP(originator messages.Address, destination messages.Address) {
	rt := a.routingTable
	destEntry, exists := rt.GetEntry(destination)
	a.seqNum++
	fmt.Println("generating route reply...")
	if !exists && destination != a.currentAddress {
		log.Printf("No route to destination address: %s\n", destination.String())
		return

	}

	if destination == a.currentAddress {
		destEntry = &RouteEntry{
			DestinationAddress: a.currentAddress,
			NextHop:            a.broadcastAddress,
			HopCount:           0,
			Precursors:         nil,
			SequenceNumber:     a.seqNum,
			Arrival:            time.Now(),
		}
	}
	rrep := &messages.RREP{
		HopCount:               destEntry.HopCount,
		DestinationAddress:     destination,
		DestinationSequenceNum: destEntry.SequenceNumber,
		OriginatorAddress:      originator,
	}
	originatorEntry, exists := rt.GetEntry(originator)

	if !exists && originator != a.broadcastAddress {
		log.Printf("No route to originator address: %s\n", originator.String())
		a.generateRRER(originator)
		return
	}
	var nextHop messages.Address
	if a.broadcastAddress == originator {
		nextHop = originator
	} else {
		nextHop = originatorEntry.NextHop
	}
	a.sendToNextHop(rrep, nextHop)
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

		a.seqNum = rrep.DestinationSequenceNum
	}

	if rrep.OriginatorAddress == a.currentAddress {
		a.dataQueue.Signal(rrep.DestinationAddress)
		return
	}

	if rrep.OriginatorAddress == a.broadcastAddress {
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
