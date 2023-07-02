package messages

import (
	"fmt"
	"strconv"
)

type RRER struct {
	DestCount                  uint8
	UnreachDestinationAddress  Address
	UnreachDestinationSequence int16
}

func (h *RRER) Type() MessageType {
	//TODO implement me
	panic("implement me")
}

func (h *RRER) HeaderSize() int {
	return 10
}

func (h *RRER) Unmarshal(data []byte) error {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return fmt.Errorf("cannot Unmarshal data: expected at least %d bytes but got %d", h.HeaderSize(), len(data))
	}

	u64, err := strconv.ParseUint(string(data[:2]), 16, 8)
	if err != nil {
		return fmt.Errorf("invalid Dest Count")
	}
	h.DestCount = uint8(u64)

	err = h.UnreachDestinationAddress.UnmarshalText(data[2:6])
	if err != nil {
		return fmt.Errorf("invalid Unreach Destination Address: %w", err)
	}

	u64, err = strconv.ParseUint(string(data[6:10]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid Unreach Dest Sequence Number")
	}
	h.UnreachDestinationSequence = int16(u64)

	return nil
}

func (h *RRER) Marshal() ([]byte, error) {
	buf := make([]byte, 0, h.HeaderSize())

	buf = append(buf, []byte(fmt.Sprintf("%02X", h.DestCount))...)
	addressBytes, err := h.UnreachDestinationAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	buf = append(buf, []byte(fmt.Sprintf("%04X", h.UnreachDestinationSequence))...)
	return buf, nil
}
