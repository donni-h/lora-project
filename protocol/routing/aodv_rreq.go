package routing

import "lora-project/protocol/messages"

func (a *AODV) handleRREQ(rreq *messages.RREQ, precursor messages.Address) {
	rt := a.routingTable
	if a.idTable.hasID(rreq.OriginatorAddress, rreq.RREQID) || rreq.OriginatorAddress == a.currentAddress {
		return
	}

	existingEntry, exists := rt.GetEntry(rreq.OriginatorAddress)

	if !exists {
		rt.AddOrUpdateEntry(
			rreq.OriginatorAddress,
			precursor,
			rreq.HopCount+1,
			[]messages.Address{},
			rreq.OriginatorSequenceNum,
		)
		a.seqNum = rreq.OriginatorSequenceNum
	} else if messages.CompareSeqnums(existingEntry.SequenceNumber, rreq.OriginatorSequenceNum) {
		rt.AddOrUpdateEntry(
			rreq.OriginatorAddress,
			precursor,
			rreq.HopCount+1,
			[]messages.Address{},
			rreq.OriginatorSequenceNum,
		)
	}
	if messages.CompareSeqnums(a.seqNum, rreq.OriginatorSequenceNum) {
		a.seqNum = rreq.OriginatorSequenceNum
	}

	a.idTable.AddID(rreq.OriginatorAddress, rreq.RREQID)

	if rreq.DestinationAddress == a.currentAddress {
		a.generateRREP(rreq.OriginatorAddress, rreq.DestinationAddress)
		return
	}

	_, exists = rt.GetEntry(rreq.DestinationAddress)

	if exists {
		a.generateRREP(rreq.OriginatorAddress, rreq.DestinationAddress)
		return
	}

	rreq.HopCount += 1
	a.sendToNextHop(rreq, a.broadcastAddress)
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
