package serial_handlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"lora-project/protocol/messages"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Command struct {
	Cmd      string
	Args     []string
	Callback func(response string, err error)
}

type ATHandler struct {
	device           io.ReadWriter
	CommandQueue     chan Command
	ErrorChan        chan error
	Done             chan bool
	responseReceived chan struct{}
	commandMutex     sync.Mutex
	commandsInFlight []Command
	MessageChan      chan messages.Message
}

func NewATHandler(device io.ReadWriter) *ATHandler {
	handler := &ATHandler{
		device:           device,
		CommandQueue:     make(chan Command),
		ErrorChan:        make(chan error),
		Done:             make(chan bool),
		responseReceived: make(chan struct{}),
		MessageChan:      make(chan messages.Message),
	}
	return handler
}

func (a *ATHandler) AddCommand(cmd Command) {
	a.CommandQueue <- cmd
}

func (a *ATHandler) Run() {
	go a.processResponses()
	a.processCommands()
}

func (a *ATHandler) processCommands() {
	for {
		select {
		case cmd := <-a.CommandQueue:
			a.commandMutex.Lock()
			a.commandsInFlight = append(a.commandsInFlight, cmd)
			err := a.sendCommand()
			if err != nil {
				a.ErrorChan <- err
				continue
			}
			time.Sleep(time.Second)
			<-a.responseReceived
		case <-a.Done:
			return
		}
	}
}

func (a *ATHandler) sendCommand() error {
	if len(a.commandsInFlight) == 0 {
		return nil
	}
	cmd := a.commandsInFlight[0]
	cmdString := cmd.Cmd
	if len(cmd.Args) > 0 {
		if cmd.Cmd == "AT+CFG" {
			cmdString += "=" + strings.Join(cmd.Args, ",")
		} else {
			cmdString += "=" + cmd.Args[0]
		}
	}
	_, err := a.device.Write([]byte(cmdString + "\r\n"))
	return err
}

func (a *ATHandler) processResponses() {
	scanner := bufio.NewScanner(a.device)
	for scanner.Scan() {
		response := strings.TrimSuffix(scanner.Text(), "\r\n")

		if strings.HasPrefix(response, "AT+SENDED") {
			a.commandMutex.Unlock()
			continue
		}

		if strings.HasPrefix(response, "LR,") {
			dataParts := strings.SplitN(response, ",", 4)
			if len(dataParts) != 4 {
				a.ErrorChan <- errors.New("Received malformed data: " + response)
				continue
			}
			sourceAddress := messages.Address{}
			err := sourceAddress.UnmarshalText([]byte(dataParts[1]))
			if err != nil {
				a.ErrorChan <- errors.New("Unable to parse source address: " + dataParts[1])
				continue
			}
			dataLength, err := strconv.ParseInt(dataParts[2], 16, 64)
			if err != nil {
				a.ErrorChan <- errors.New("Unable to parse data length: " + dataParts[2])
				continue
			}
			payload := dataParts[3]
			if int64(len(payload)) != dataLength {
				a.ErrorChan <- errors.New(fmt.Sprintf("Data length mismatch: expected %d, got %d", dataLength, len(payload)))
				continue
			}
			message, err := messages.Unmarshal([]byte(payload))
			if err != nil {
				a.ErrorChan <- err
				continue
			}
			// handle the received data
			fmt.Println("Received data from", sourceAddress.String(), ": message type ", message.Type())
			a.MessageChan <- message
			continue
		}

		a.commandMutex.Lock()
		if len(a.commandsInFlight) == 0 {
			a.commandMutex.Unlock()
			continue
		}
		cmd := a.commandsInFlight[0]
		a.commandsInFlight = a.commandsInFlight[1:]

		if strings.HasPrefix(response, "AT+ERR") {
			cmd.Callback("", errors.New(response))
		} else if strings.HasPrefix(response, "AT+OK") {
			cmd.Callback(response, nil)
		}
		a.commandMutex.Unlock()
		a.responseReceived <- struct{}{}
	}
}
