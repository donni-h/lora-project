package routing

import "time"

type RouteEntry struct {
	DestinationAddress     [4]byte
	NextHop                [4]byte
	HopCount               [2]byte
	DestinationSequenceNum [4]byte
	Arrival                time.Time
}
