package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RREP struct {
	T                      byte
	HopCount               [2]byte
	DestinationAddress     [4]byte
	DestinationSequenceNum [4]byte
	OriginatorAddress      [4]byte
}

func (h *RREP) HeaderSize() int {
	return 15
}

func (h *RREP) unmarshal(data []byte) error {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return fmt.Errorf("cannot unmarshal data: expected %d bytes but got %d", h.HeaderSize(), len(data))
	}

	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.BigEndian, h)
	if err != nil {
		return err
	}

	return nil
}

func (h *RREP) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, h.HeaderSize()))

	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
