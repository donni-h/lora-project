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

func (r *RRER) Type() MessageType {
	return TypeRRER
}

func (r *RRER) HeaderSize() int {
	return 5
}

func (r *RRER) Unmarshal(data []byte) error {
	if len(data) < 10 {
		// Handle error: insufficient data
		return fmt.Errorf("cannot Unmarshal data: expected at least %d bytes but got %d", r.HeaderSize(), len(data))
	}

	u64, err := strconv.ParseUint(string(data[:2]), 16, 8)
	if err != nil {
		return fmt.Errorf("invalid Dest Count")
	}
	r.DestCount = uint8(u64)

	err = r.UnreachDestinationAddress.UnmarshalText(data[2:6])
	if err != nil {
		return fmt.Errorf("invalid Unreach Destination Address: %w", err)
	}

	u64, err = strconv.ParseUint(string(data[6:10]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid Unreach Dest Sequence Number")
	}
	r.UnreachDestinationSequence = int16(u64)

	return nil
}

func (r *RRER) Marshal() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, byte(r.Type()))
	buf = append(buf, []byte(fmt.Sprintf("%02X", r.DestCount))...)
	addressBytes, err := r.UnreachDestinationAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	buf = append(buf, []byte(fmt.Sprintf("%04X", uint16(r.UnreachDestinationSequence)))...)
	return buf, nil
}
