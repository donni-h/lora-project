package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RRER struct {
	T                          byte
	DestCount                  [2]byte
	UnreachDestinationAddress  [4]byte
	UnreachDestinationSequence [4]byte
}

func (h *RRER) HeaderSize() int {
	return 11
}

func (h *RRER) unmarshal(data []byte) error {
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

func (h *RRER) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, h.HeaderSize()))

	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
