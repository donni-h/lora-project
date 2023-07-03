package messages

import (
	"fmt"
	"strconv"
)

type Data struct {
	DestinationAddress Address
	OriginatorAddress  Address
	DataSequenceNumber int16
	Payload            []byte
}

func (r *Data) Type() MessageType {
	return Type_Data
}

func (r *Data) HeaderSize() int {
	return 6
}

func (r *Data) Unmarshal(data []byte) error {
	if len(data) < 12 {
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
	u64, err := strconv.ParseUint(string(data[8:12]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid Unreach Dest Sequence Number")
	}
	r.DataSequenceNumber = int16(u64)

	r.Payload = data[12:]
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
