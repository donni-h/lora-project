package serial_handlers

import (
	"bytes"
	"io"
	"log"
	"lora-project/protocol/messages"
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
		// Write a valid Data message to the Pipe
		_, err := w.Write([]byte("LR,4761,12,327BCD1230001HELLO\r\n"))
		if err != nil {
			atHandler.ErrorChan <- err
		}
		err = w.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Listen for the message
	select {
	case err := <-atHandler.ErrorChan:
		t.Fatal(err)

	case msg := <-atHandler.MessageChan:
		// Here, we can test if the message is processed correctly
		dataMsg, ok := msg.(*messages.Data)
		if !ok {
			t.Fatal("Expected a Data message")
		}

		// Assert various properties of the message
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
		t.Fatal("Timeout waiting for message")
	}
}
