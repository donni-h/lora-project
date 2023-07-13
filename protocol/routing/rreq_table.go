package routing

import (
	"lora-project/protocol/messages"
	"sync"
)

type RREQIDTable struct {
	ids map[messages.Address][]uint16
	sync.RWMutex
}

func NewRREQIDTable() *RREQIDTable {
	return &RREQIDTable{
		ids: make(map[messages.Address][]uint16),
	}
}

func (rt *RREQIDTable) AddID(address messages.Address, id uint16) {
	rt.Lock()
	defer rt.Unlock()
	rt.ids[address] = append(rt.ids[address], id)
}

func (rt *RREQIDTable) hasID(address messages.Address, id uint16) bool {
	rt.RLock()
	defer rt.RUnlock()
	for _, v := range rt.ids[address] {
		if v == id {
			return true
		}
	}
	return false

}
