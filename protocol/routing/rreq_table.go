package routing

import "lora-project/protocol/messages"

type RREQIDTable map[messages.Address][]uint16

func NewRREQIDTable() RREQIDTable {
	return make(RREQIDTable)
}

func (rt RREQIDTable) AddID(address messages.Address, id uint16) {
	rt[address] = append(rt[address], id)
}

func (rt RREQIDTable) hasID(address messages.Address, id uint16) bool {
	for _, v := range rt[address] {
		if v == id {
			return true
		}
	}
	return false

}
