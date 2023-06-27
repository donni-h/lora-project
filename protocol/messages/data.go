package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Data struct {
	T                  byte
	DestinationAddress [4]byte
	OriginatorAddress  [4]byte
	DataSequenceNumber [2]byte
	Payload            []byte
}

func (h *Data) HeaderSize() int {
	return 11
}

func (h *Data) unmarshal(data []byte) error {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return fmt.Errorf("cannot unmarshal data: expected at least %d bytes but got %d", h.HeaderSize(), len(data))
	}

	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.BigEndian, h)
	if err != nil {
		return err
	}

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
