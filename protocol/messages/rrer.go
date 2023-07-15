package messages

import (
	"fmt"
)

type RRER struct {
	UnreachDestinationAddress Address
}

func (r *RRER) Type() MessageType {
	return TypeRRER
}

func (r *RRER) HeaderSize() int {
	return 4
}

func (r *RRER) Unmarshal(data []byte) error {
	if len(data) < 4 {
		// Handle error: insufficient data
		return fmt.Errorf("cannot Unmarshal data: expected at least %d bytes but got %d", r.HeaderSize(), len(data))
	}

	err := r.UnreachDestinationAddress.UnmarshalText(data[:4])
	if err != nil {
		return fmt.Errorf("invalid Unreach Destination Address: %w", err)
	}

	return nil
}

func (r *RRER) Marshal() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, byte(r.Type()))
	addressBytes, err := r.UnreachDestinationAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	return buf, nil
}
