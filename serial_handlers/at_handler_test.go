package serial_handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"lora-project/protocol/messages"
	"strconv"
	"strings"
	"testing"
	"time"
)

type ReadWriter struct {
	*io.PipeReader
	*io.PipeWriter
}

func NewReadWriter(r *io.PipeReader, w *io.PipeWriter) *ReadWriter {
	return &ReadWriter{r, w}
}

func TestConfigTransceiver(t *testing.T) {
	buffer := &bytes.Buffer{}
	atHandler := NewATHandler(buffer)

	go func() {
		atHandler.Configure([]string{"433000000", "20", "6", "10", "1", "1", "0", "0", "0", "0", "3000", "8", "4"})
	}()
	select {
	case cmd := <-atHandler.CommandQueue:
		cmdString := cmd.Cmd + "=" + strings.Join(cmd.Args, ",")
		expectedCmd := "AT+CFG=433000000,20,6,10,1,1,0,0,0,0,3000,8,4"
		if cmdString != expectedCmd {
			t.Errorf("Unexpected command: %v", cmdString)
		}
	}
}

func TestProcessIncomingData(t *testing.T) {
	r, w := io.Pipe()
	rw := NewReadWriter(r, w)
	atHandler := NewATHandler(rw)
	go atHandler.Run()
	go func() {
		// Write a valid Data Message to the Pipe
		_, err := w.Write([]byte("LR,4761,12,327BCD1230001HELLO\r\n"))
		if err != nil {
			atHandler.ErrorChan <- err
		}
		err = w.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Listen for the Message
	select {
	case err := <-atHandler.ErrorChan:
		t.Fatal(err)

	case msgEvent := <-atHandler.MessageChan:
		msg := msgEvent.Message
		// Here, we can test if the Message is processed correctly
		dataMsg, ok := msg.(*messages.Data)
		if !ok {
			t.Fatal("Expected a Data Message")
		}

		// Assert various properties of the Message
		expectedDestination := messages.Address{}
		_ = expectedDestination.UnmarshalText([]byte("27BC"))
		if dataMsg.DestinationAddress != expectedDestination {
			t.Fatalf("Expected DestinationAddress to be %v but got %v", expectedDestination, dataMsg.DestinationAddress)
		}
		expectedOrigin := messages.Address{}
		_ = expectedOrigin.UnmarshalText([]byte("D123"))
		if dataMsg.OriginatorAddress != expectedOrigin {
			t.Fatalf("Expected OriginatorAddress to be %v but got %v", expectedOrigin, dataMsg.OriginatorAddress)
		}
		if dataMsg.DataSequenceNumber != 1 {
			t.Fatalf("Expected DataSequenceNumber to be 1 but got %d", dataMsg.DataSequenceNumber)
		}
		expectedPayload := []byte("HELLO")
		if string(dataMsg.Payload) != string(expectedPayload) {
			t.Fatalf("Expected Payload to be %s but got %s", string(expectedPayload), string(dataMsg.Payload))
		}

	case <-time.After(100 * time.Second):
		t.Fatal("Timeout waiting for Message")
	}
}

type MockReadWriter struct {
	Buffer     bytes.Buffer
	WriteCount int
}

func (rw *MockReadWriter) GetData() string {
	return rw.Buffer.String()
}

func (rw *MockReadWriter) Write(p []byte) (n int, err error) {
	rw.WriteCount++
	fmt.Println("wrote.")
	_, err = rw.Buffer.Write(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (rw *MockReadWriter) Read(p []byte) (n int, err error) {
	return rw.Buffer.Read(p)
}

func TestSendingData(t *testing.T) {
	device := &MockReadWriter{}
	handler := NewATHandler(device)
	message := []byte{'0', '2', '1', '1', '1', '1', '0', '0', '0', '0', '4', '7', 'H', 'A', 'L', 'L', 'O', 'W', 'E', 'L', 'T', '!'}
	var msg messages.Data
	err := msg.Unmarshal(message)
	if err != nil {
		return
	}
	handler.SendMessage(&msg)
	select {
	case data := <-handler.MessageChan:
		fmt.Printf("%+v\n", data)
	case err = <-handler.ErrorChan:
		fmt.Println(err)
	case cmd := <-handler.CommandQueue:
		if cmd.Cmd != "AT+SEND" {
			t.Fatalf("Excepted command to be AT+SEND but was: %s", cmd.Cmd)
		}

		if cmd.Args[0] != strconv.Itoa(len(message)+1) {
			fmt.Println(string(message))
			fmt.Println(cmd)
			t.Fatalf("Excepted Data Length to be %v but got: %s", len(message)+1, cmd.Args[0])
		}
		cmd.Callback("AT+OK", nil)
		scanner := bufio.NewScanner(device)
		scanner.Scan()
		var resp messages.Data
		err = resp.Unmarshal([]byte(scanner.Text()))
		if err != nil {
			t.Fatalf("Couldnt unmarshal data: %s", err)
		}
		if string(resp.Payload) != "7HALLOWELT!" {
			t.Fatalf("Wrong payload. %s", string(resp.Payload))
		}
		fmt.Println(string(resp.Payload))
	case <-time.After(time.Second * 50):
		t.Fatalf("Nothing happened")
	}
	fmt.Println(device.WriteCount)
}
