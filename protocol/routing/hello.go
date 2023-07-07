package routing

import (
	"lora-project/protocol/messages"
	"time"
)

const (
	HelloInterval = time.Minute
)

func (a *AODV) helloRoutine() {
	ticker := time.NewTicker(HelloInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.sendHello()
		}
	}
}

func (a *AODV) sendHello() {
	rrep := &messages.RREP{
		HopCount:               0,
		DestinationAddress:     a.currentAddress,
		DestinationSequenceNum: a.seqNum,
		OriginatorAddress:      a.currentAddress,
	}

	a.handler.SendMessage(rrep)
}
