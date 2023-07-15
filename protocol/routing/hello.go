package routing

import (
	"time"
)

const (
	HelloInterval = time.Second * 10
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

	a.generateRREP(a.broadcastAddress, a.currentAddress)
}
