package serial_handlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"lora-project/protocol/messages"
	"strconv"
	"strings"
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
	currentCommand   *Command
	MessageChan      chan MessageEvent
}

func NewATHandler(device io.ReadWriter) *ATHandler {
	handler := &ATHandler{
		device:           device,
		CommandQueue:     make(chan Command, 10),
		ErrorChan:        make(chan error),
		Done:             make(chan bool),
		responseReceived: make(chan struct{}),
		MessageChan:      make(chan MessageEvent, 10),
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
	fmt.Println(string(data))
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

func (a *ATHandler) SetOwnAddress(addr messages.Address) {
	cmd := Command{
		Cmd:  "AT+ADDR",
		Args: []string{addr.String()},
		Callback: func(response string, err error) {
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
func (a *ATHandler) SetReceive() {
	cmd := Command{
		Cmd:  "AT+RX",
		Args: nil,
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
			a.currentCommand = &cmd
			err := a.sendCommand()
			if err != nil {
				a.ErrorChan <- err
				continue
			}
			<-a.responseReceived
		case <-a.Done:
			return
		}
	}
}

func (a *ATHandler) sendCommand() error {

	cmd := a.currentCommand
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
	reader := bufio.NewReader(a.device)
	for {
		responseType, err := reader.ReadString(',')
		if err != nil {
			if err == io.EOF {
				break
			}
			a.ErrorChan <- err
			return
		}
		responseType = strings.TrimSpace(responseType)
		if strings.HasPrefix(responseType, "LR,") {
			a.handleReceivedData(reader)
			continue
		}
		response, err := reader.ReadString('\n')
		fmt.Println(response)
		response = fmt.Sprintf("%s%s", responseType, strings.TrimSpace(response))
		if a.currentCommand.Cmd == "AT+SEND" {
			a.handleCommandSent(response)
		} else {
			a.handleCommandResponse(response)
		}
	}
}

func (a *ATHandler) handleReceivedData(reader *bufio.Reader) {

	srcBytes, err := reader.ReadBytes(',')
	if err != nil {
		a.ErrorChan <- err
		return
	}
	srcBytes = srcBytes[:len(srcBytes)-1]
	srcAddress := messages.Address{}

	err = srcAddress.UnmarshalText(srcBytes)
	if err != nil {
		a.ErrorChan <- err
		return
	}
	lengthStr, err := reader.ReadString(',')
	if err != nil {
		a.ErrorChan <- err
		return
	}
	lengthStr = strings.TrimSpace(lengthStr)
	lengthStr = strings.TrimSuffix(lengthStr, ",")
	length, err := strconv.ParseInt(lengthStr, 16, 64) // convert from hex string to int
	if err != nil {
		a.ErrorChan <- err
		return
	}

	payload := make([]byte, length)

	_, err = io.ReadFull(reader, payload)
	if err != nil {
		a.ErrorChan <- err
		return
	}
	fmt.Println(string(payload), srcAddress.String())
	msg, err := messages.Unmarshal(payload)
	if err != nil {
		a.ErrorChan <- err
		return
	}
	event := MessageEvent{
		Message:   msg,
		Precursor: srcAddress,
	}
	if srcAddress.String() != "7890" {
		return
	}

	a.MessageChan <- event
}

func (a *ATHandler) handleCommandSent(response string) {
	if response == "AT,OK" {
		a.currentCommand.Callback(response, nil)
	} else if response == "AT,SENDING" {
	} else if response == "AT,SENDED" {
		a.currentCommand = nil
		a.responseReceived <- struct{}{}
	} else {
		a.ErrorChan <- errors.New("unexpected response received")
	}
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
	if a.currentCommand == nil {
		a.ErrorChan <- errors.New("unexpected response: " + response)
		return
	}
	cmd := a.currentCommand
	a.currentCommand = nil
	cmd.Callback(response, nil)
	a.responseReceived <- struct{}{}
}
