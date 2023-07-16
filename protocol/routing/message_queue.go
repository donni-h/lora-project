package routing

import (
	"errors"
	"lora-project/protocol/messages"
	"sync"
	"time"
)

var ErrTimeOut = errors.New("RREQ timeout")

type DataQueue struct {
	messages map[messages.Address][]*messages.Data
	signals  map[messages.Address]chan struct{}
	mux      sync.Mutex
}

func NewDataQueue() *DataQueue {
	return &DataQueue{
		messages: make(map[messages.Address][]*messages.Data),
		signals:  make(map[messages.Address]chan struct{}),
	}
}

func (q *DataQueue) Push(data *messages.Data) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.messages[data.DestinationAddress] = append(q.messages[data.DestinationAddress], data)
	if _, ok := q.signals[data.DestinationAddress]; !ok {
		q.signals[data.DestinationAddress] = make(chan struct{})
	}
}

func (q *DataQueue) Pop(destination messages.Address) ([]*messages.Data, error) {
	q.mux.Lock()
	signal, ok := q.signals[destination]
	q.mux.Unlock()
	if !ok {
		return nil, nil
	}

	timeoutChan := time.After(time.Second * 10)

	select {
	case <-signal:
	// continue as usual
	case <-timeoutChan:
		return nil, ErrTimeOut

	}

	q.mux.Lock()
	defer q.mux.Unlock()

	msgQueue, ok := q.messages[destination]

	if !ok || len(msgQueue) == 0 {
		return nil, nil
	}

	delete(q.messages, destination)
	delete(q.signals, destination)
	return msgQueue, nil

}

func (q *DataQueue) Signal(destination messages.Address) {
	q.mux.Lock()
	signal, ok := q.signals[destination]
	q.mux.Unlock()
	if !ok {
		return
	}
	close(signal)
}
