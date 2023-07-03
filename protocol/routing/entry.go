package routing

import (
	"lora-project/protocol/messages"
	"time"
)

type RouteEntry struct {
	DestinationAddress     messages.Address
	NextHop                messages.Address
	HopCount               uint8
	Precursors             []messages.Address
	DestinationSequenceNum int16
	Arrival                time.Time
}
