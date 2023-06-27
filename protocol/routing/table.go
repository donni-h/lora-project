package routing

import "time"

type Table map[[4]byte]*RouteEntry

func NewTable() Table {
	return make(Table)
}

func (rt Table) AddOrUpdateEntry(address [4]byte, nextHop [4]byte, hopCount [2]byte, seqNum [4]byte) {
	rt[address] = &RouteEntry{
		NextHop:                nextHop,
		HopCount:               hopCount,
		DestinationSequenceNum: seqNum,
		Arrival:                time.Now(),
	}
}

func (rt Table) GetEntry(address [4]byte) (*RouteEntry, bool) {
	entry, found := rt[address]
	return entry, found
}

func (rt Table) DeleteEntry(address [4]byte) {
	delete(rt, address)
}
