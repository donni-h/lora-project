package messages

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

type RREP struct {
	HopCount               uint8
	DestinationAddress     Address
	DestinationSequenceNum int16
	OriginatorAddress      Address
}

var _ Message = &RREP{}

func (r RREP) Type() MessageType {
	return TypeRREP
}

func (r RREP) HeaderSize() int {
	return 7
}

func (r *RREP) Unmarshal(data []byte) error {
	if len(data) < 14 {
		return fmt.Errorf("wrong data size")
	}

	u64, err := strconv.ParseUint(string(data[:2]), 16, 8)
	if err != nil {
		return err
	}

	r.HopCount = uint8(u64)

	err = r.DestinationAddress.UnmarshalText(data[2:6])
	if err != nil {
		return err
	}

	bytes, err := hex.DecodeString(string(data[6:10]))
	if err != nil {
		return err
	}
	r.DestinationSequenceNum = int16(binary.BigEndian.Uint16(bytes))

	err = r.OriginatorAddress.UnmarshalText(data[10:14])
	if err != nil {
		return err
	}

	return nil
}

func (r RREP) Marshal() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, byte(r.Type()))
	buf = append(buf, []byte(fmt.Sprintf("%02X", r.HopCount))...)
	addressBytes, err := r.DestinationAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	buf = append(buf, []byte(fmt.Sprintf("%04X", r.DestinationSequenceNum))...)
	addressBytes, err = r.OriginatorAddress.MarshalText()
	if err != nil {
		return nil, err
	}
	buf = append(buf, addressBytes...)
	return buf, nil
}
