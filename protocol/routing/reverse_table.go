package routing

import "time"

type ReverseTable map[[4]byte]*RouteEntry

func NewReverseRoutingTable() ReverseTable {
	return make(ReverseTable)
}

func (rrt ReverseTable) AddOrUpdateEntry(address [4]byte, nextHop [4]byte, hopCount [2]byte, seqNum [4]byte) {
	rrt[address] = &RouteEntry{
		NextHop:                nextHop,
		HopCount:               hopCount,
		DestinationSequenceNum: seqNum,
		Arrival:                time.Now(),
	}
}
