package routing

import (
	"lora-project/protocol/messages"
	"time"
)

type Table map[messages.Address]*RouteEntry

func NewTable() Table {
	return make(Table)
}

func (rt Table) AddOrUpdateEntry(address messages.Address, nextHop messages.Address, hopCount uint8, precursors []messages.Address, seqNum int16) {
	rt[address] = &RouteEntry{
		NextHop:        nextHop,
		HopCount:       hopCount,
		Precursors:     precursors,
		SequenceNumber: seqNum,
		Arrival:        time.Now(),
	}
}

func (rt Table) GetEntry(address messages.Address) (*RouteEntry, bool) {
	entry, found := rt[address]
	return entry, found
}

func (rt Table) DeleteEntry(address messages.Address) {
	delete(rt, address)
}

func (rt Table) AddPrecursor(address messages.Address, precursor messages.Address) bool {
	entry, exists := rt[address]
	if !exists {
		return false
	}
	entry.Precursors = append(entry.Precursors, precursor)
	return true
}
