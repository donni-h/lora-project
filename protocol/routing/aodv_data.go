package routing

import (
	"fmt"
	"log"
	"lora-project/protocol/messages"
)

func (a *AODV) SendData(payload string, destination messages.Address) {
	a.seqNum++
	data := &messages.Data{
		DestinationAddress: destination,
		OriginatorAddress:  a.currentAddress,
		Payload:            []byte(payload),
	}

	entry, exists := a.routingTable.GetEntry(destination)
	if !exists {
		_, exists = a.dataQueue.signals[destination]
		a.dataQueue.Push(data)

		if !exists {
			a.generateRREQ(destination)
			go a.queueData(destination)
		}
		return
	}
	a.sendToNextHop(data, entry.NextHop)
}

func (a *AODV) queueData(destination messages.Address) {
	timeOutCount := 0

	for {
		pendingData, err := a.dataQueue.Pop(destination)
		if err == ErrTimeOut {
			timeOutCount++
			if timeOutCount == 3 {
				a.dataQueue.mux.Lock()
				delete(a.dataQueue.messages, destination)
				delete(a.dataQueue.signals, destination)
				a.dataQueue.mux.Unlock()
				return
			}
			a.generateRREQ(destination)
			continue
		}

		fmt.Printf("das sind die queued messages: %+v\n", pendingData)
		for _, msg := range pendingData {
			entry, ok := a.routingTable.GetEntry(destination)

			if !ok {
				log.Println("Couldn't send message, no route to:" + destination.String())
				return
			}

			a.sendToNextHop(msg, entry.NextHop)
		}
		return
	}
}

func (a *AODV) handleData(data *messages.Data) {

	if data.DestinationAddress == a.currentAddress {
		a.IncomingDataQueue <- *data
		return
	}

	route, exists := a.routingTable.GetEntry(data.DestinationAddress)
	if !exists {
		a.generateRRER(data.DestinationAddress)
		return
	}
	a.sendToNextHop(data, route.NextHop)
}
