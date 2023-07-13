package routing

import (
	"lora-project/protocol/messages"
	"sync"
)

type DataQueue struct {
	messages map[messages.Address][]*messages.Data
	conds    map[messages.Address]*sync.Cond
	mux      sync.Mutex
}

func NewDataQueue() *DataQueue {
	return &DataQueue{
		messages: make(map[messages.Address][]*messages.Data),
		conds:    make(map[messages.Address]*sync.Cond),
	}
}

func (q *DataQueue) Push(data *messages.Data) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.messages[data.DestinationAddress] = append(q.messages[data.DestinationAddress], data)
	if _, ok := q.conds[data.DestinationAddress]; !ok {
		q.conds[data.DestinationAddress] = sync.NewCond(&sync.Mutex{})
	}
}

func (q *DataQueue) Pop(destination messages.Address) []*messages.Data {
	q.mux.Lock()
	cond, ok := q.conds[destination]
	q.mux.Unlock()
	if !ok {
		return nil
	}

	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()

	q.mux.Lock()
	defer q.mux.Unlock()

	msgQueue, ok := q.messages[destination]

	if !ok || len(msgQueue) == 0 {
		return nil
	}

	delete(q.messages, destination)
	delete(q.conds, destination)

	return msgQueue

}

func (q *DataQueue) Signal(destination messages.Address) {
	q.mux.Lock()
	cond, ok := q.conds[destination]
	q.mux.Unlock()
	if !ok {
		return
	}
	cond.Broadcast()
}
