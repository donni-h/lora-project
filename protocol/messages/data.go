package messages

import (
	"fmt"
)

type Data struct {
	DestinationAddress Address
	OriginatorAddress  Address
	Payload            []byte
}

func (r *Data) Type() MessageType {
	return TypeData
}

func (r *Data) HeaderSize() int {
	return 8
}

func (r *Data) Unmarshal(data []byte) error {
	if len(data) < 8 {
		// Handle error: insufficient data
		return fmt.Errorf("cannot Unmarshal data: expected at least %d bytes but got %d", 10, len(data))
	}

	err := r.DestinationAddress.UnmarshalText(data[:4])
	if err != nil {
		return fmt.Errorf("invalid Destination Address")
	}

	err = r.OriginatorAddress.UnmarshalText(data[4:8])
	if err != nil {
		return fmt.Errorf("invalid Originator Address")
	}

	r.Payload = data[8:]
	return nil
}

func (r *Data) Marshal() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, byte(r.Type()))
	addressBytes, err := r.DestinationAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	addressBytes, err = r.OriginatorAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	buf = append(buf, r.Payload...)

	return buf, nil
}
