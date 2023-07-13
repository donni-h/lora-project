package routing

import (
	"lora-project/protocol/messages"
	"sync"
	"time"
)

type Table struct {
	sync.RWMutex
	routes map[messages.Address]*RouteEntry
}

func NewTable() *Table {
	return &Table{
		routes: make(map[messages.Address]*RouteEntry),
	}
}

func (rt *Table) AddOrUpdateEntry(address messages.Address, nextHop messages.Address, hopCount uint8, precursors []messages.Address, seqNum int16) {
	rt.Lock()
	defer rt.Unlock()
	rt.routes[address] = &RouteEntry{
		NextHop:        nextHop,
		HopCount:       hopCount,
		Precursors:     precursors,
		SequenceNumber: seqNum,
		Arrival:        time.Now(),
	}
}

func (rt *Table) GetEntry(address messages.Address) (*RouteEntry, bool) {
	rt.RLock()
	defer rt.RUnlock()
	entry, found := rt.routes[address]
	return entry, found
}

func (rt *Table) DeleteEntry(address messages.Address) {
	rt.Lock()
	defer rt.Unlock()
	delete(rt.routes, address)
}

func (rt *Table) AddPrecursor(address messages.Address, precursor messages.Address) bool {
	rt.Lock()
	defer rt.Unlock()
	entry, exists := rt.routes[address]
	if !exists {
		return false
	}
	entry.Precursors = append(entry.Precursors, precursor)
	return true
}
