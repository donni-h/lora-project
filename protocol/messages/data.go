package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Data struct {
	DestinationAddress [4]byte
	OriginatorAddress  [4]byte
	DataSequenceNumber [2]byte
	Payload            []byte
}

func (h *Data) Type() MessageType {
	//TODO implement me
	panic("implement me")
}

func (h *Data) HeaderSize() int {
	return 10
}

func (h *Data) Unmarshal(data []byte) error {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return fmt.Errorf("cannot Unmarshal data: expected at least %d bytes but got %d", h.HeaderSize(), len(data))
	}

	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.BigEndian, &h.DestinationAddress)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &h.OriginatorAddress)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &h.DataSequenceNumber)
	if err != nil {
		return err
	}

	// Read remaining data into Payload
	h.Payload = buf.Bytes()

	return nil
}

func (h *Data) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, h.HeaderSize()))

	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
