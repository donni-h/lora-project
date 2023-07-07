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

type MessageEvent struct {
	Message   messages.Message
	Precursor messages.Address
}

type ATHandler struct {
	device           io.ReadWriter
	CommandQueue     chan Command
	ErrorChan        chan error
	Done             chan bool
	responseReceived chan struct{}
	commandMutex     sync.Mutex
	commandsInFlight []Command
	MessageChan      chan MessageEvent
}

func NewATHandler(device io.ReadWriter) *ATHandler {
	handler := &ATHandler{
		device:           device,
		CommandQueue:     make(chan Command, 10),
		ErrorChan:        make(chan error),
		Done:             make(chan bool),
		responseReceived: make(chan struct{}),
		MessageChan:      make(chan MessageEvent),
	}
	go handler.Run()
	return handler
}

func (a *ATHandler) AddCommand(cmd Command) {
	a.CommandQueue <- cmd
}

// SendMessage sends a messages.Message via the ATHandler.
func (a *ATHandler) SendMessage(msg messages.Message) {
	data, err := msg.Marshal()
	if err != nil {
		a.ErrorChan <- err
		return
	}
	cmd := Command{
		Cmd:  "AT+SEND",
		Args: []string{strconv.Itoa(len(data))},
		Callback: func(response string, err error) {
			if err != nil {
				a.ErrorChan <- err
				return
			}
			_, err = a.device.Write(data)
			if err != nil {
				a.ErrorChan <- err
			}
		},
	}
	a.AddCommand(cmd)
}

// SetTargetAddress sets the destination address for the ATHandler.
func (a *ATHandler) SetTargetAddress(addr messages.Address) {
	cmd := Command{
		Cmd:  "AT+DEST",
		Args: []string{addr.String()},
		Callback: func(response string, err error) {
			if err != nil {
				a.ErrorChan <- err
			}
		},
	}
	a.AddCommand(cmd)
}

// Configure configures the LoRa transceiver using the AT+CFG command.
func (a *ATHandler) Configure(args []string) {
	cmd := Command{
		Cmd:  "AT+CFG",
		Args: args,
		Callback: func(response string, err error) {
			if err != nil {
				a.ErrorChan <- err
			}
		},
	}
	a.AddCommand(cmd)
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
			a.handleCommandSent()
		} else if strings.HasPrefix(response, "LR,") {
			a.handleReceivedData(response)
		} else {
			a.handleCommandResponse(response)
		}
	}
}

func (a *ATHandler) handleReceivedData(response string) {
	parts := strings.Split(response, ",")
	if len(parts) != 4 {
		a.ErrorChan <- errors.New("malformed received data")
		return
	}

	srcAddress := messages.Address{}
	err := srcAddress.UnmarshalText([]byte(parts[1]))
	if err != nil {
		a.ErrorChan <- err
		return
	}

	dataLen, err := strconv.ParseInt(parts[2], 16, 32)
	if err != nil {
		a.ErrorChan <- err
		return
	}

	data := parts[3]
	if len(data) != int(dataLen) {
		a.ErrorChan <- errors.New("received data length does not match expected length")
		return
	}

	msg, err := messages.Unmarshal([]byte(data))
	if err != nil {
		a.ErrorChan <- err
		return
	}
	event := MessageEvent{
		Message:   msg,
		Precursor: srcAddress,
	}
	a.MessageChan <- event
}

func (a *ATHandler) handleCommandSent() {
	a.commandMutex.Lock()
	defer a.commandMutex.Unlock()

	if len(a.commandsInFlight) == 0 {
		a.ErrorChan <- errors.New("unexpected 'AT+SENDED'")
		return
	}
	cmd := a.commandsInFlight[0]
	a.commandsInFlight = a.commandsInFlight[1:]
	cmd.Callback("AT+SENDED", nil)
}

func (a *ATHandler) SendData(msg messages.Message) error {
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	dataLen := len(data)
	if dataLen > 250 {
		return fmt.Errorf("data length exceeds the maximum limit of 250 bytes")
	}

	a.AddCommand(Command{
		Cmd:  "AT+SEND",
		Args: []string{strconv.Itoa(dataLen)},
		Callback: func(response string, err error) {
			if err != nil {
				a.ErrorChan <- err
				return
			}
			if response != "AT,OK" {
				a.ErrorChan <- errors.New("failed to send data: " + response)
				return
			}
			_, err = a.device.Write(data)
			if err != nil {
				a.ErrorChan <- err
				return
			}
		},
	})
	return nil
}

func (a *ATHandler) handleCommandResponse(response string) {
	a.commandMutex.Lock()
	defer a.commandMutex.Unlock()

	if len(a.commandsInFlight) == 0 {
		a.ErrorChan <- errors.New("unexpected response: " + response)
		return
	}
	cmd := a.commandsInFlight[0]
	a.commandsInFlight = a.commandsInFlight[1:]
	cmd.Callback(response, nil)
}
