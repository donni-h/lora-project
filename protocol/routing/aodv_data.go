package routing

import (
	"log"
	"lora-project/protocol/messages"
)

func (a *AODV) SendData(payload string, destination messages.Address) {
	a.seqNum++
	data := &messages.Data{
		DestinationAddress: destination,
		OriginatorAddress:  a.currentAddress,
		DataSequenceNumber: a.seqNum,
		Payload:            []byte(payload),
	}

	entry, exists := a.routingTable.GetEntry(destination)
	if !exists {
		_, exists = a.dataQueue.conds[destination]
		a.dataQueue.Push(data)

		if !exists {
			a.generateRREQ(destination)
			go a.dataQueue.Pop(destination)
		}
		return
	}
	a.sendToNextHop(data, entry.NextHop)
}

func (a *AODV) queueData(destination messages.Address) {
	go func() {
		pendingData := a.dataQueue.Pop(destination)

		for _, msg := range pendingData {
			entry, ok := a.routingTable.GetEntry(destination)

			if !ok {
				log.Println("Couldn't send message, no route to:" + destination.String())
				return
			}

			a.sendToNextHop(msg, entry.NextHop)
		}
	}()
}

func (a *AODV) handleData(data *messages.Data) {
	if messages.CompareSeqnums(a.seqNum, data.DataSequenceNumber) {
		a.seqNum = data.DataSequenceNumber
	}
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
