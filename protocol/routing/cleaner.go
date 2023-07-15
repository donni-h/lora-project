package routing

import (
	"fmt"
	"lora-project/protocol/messages"
	"time"
)

const (
	invalidTimeframe = time.Second * 25
)

func (a *AODV) StartExpirationWorker() {
	go func() {
		for {
			<-time.After(invalidTimeframe)
			now := time.Now()

			var toDelete []messages.Address
			a.routingTable.RLock()
			for address, entry := range a.routingTable.routes {
				if now.Sub(entry.Arrival) > invalidTimeframe {
					toDelete = append(toDelete, address)
				}
			}
			a.routingTable.RUnlock()
			// perform deletion
			for _, address := range toDelete {
				fmt.Println("Deleting address:" + address.String())
				a.routingTable.DeleteEntry(address)
				a.generateRRER(address)
			}
		}
	}()
}
