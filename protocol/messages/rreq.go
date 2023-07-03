package messages

import (
	"fmt"
	"strconv"
)

type RREQ struct {
	UFlag                  bool
	HopCount               uint8
	RREQID                 uint16
	DestinationAddress     Address
	DestinationSequenceNum int16
	OriginatorAddress      Address
	OriginatorSequenceNum  int16
}

func (r *RREQ) Type() MessageType {
	return Type_RREQ
}

func (r *RREQ) HeaderSize() int {
	return 18
}

func (r *RREQ) Unmarshal(data []byte) error {
	if len(data) < 23 {
		return fmt.Errorf("cannot Unmarshal data: expected %d bytes but got %d", 23, len(data))
	}

	r.UFlag = data[0] == 'Y'

	var err error
	u64, err := strconv.ParseUint(string(data[1:3]), 16, 8)
	if err != nil {
		return err
	}
	r.HopCount = uint8(u64)

	u64, err = strconv.ParseUint(string(data[3:7]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid RREQ ID")
	}
	r.RREQID = uint16(u64)

	err = r.DestinationAddress.UnmarshalText(data[7:11])
	if err != nil {
		return fmt.Errorf("invalid Destination Address")
	}

	u64, err = strconv.ParseUint(string(data[11:15]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid Destination Sequence Number")
	}
	r.DestinationSequenceNum = int16(u64)

	err = r.OriginatorAddress.UnmarshalText(data[15:19])
	if err != nil {
		return fmt.Errorf("invalid Originator Address")
	}

	u64, err = strconv.ParseUint(string(data[19:23]), 16, 16)
	if err != nil {
		return fmt.Errorf("invalid Originator Sequence Number")
	}
	r.OriginatorSequenceNum = int16(u64)

	return nil
}

func (r *RREQ) Marshal() ([]byte, error) {
	buf := []byte{}
	buf = append(buf, byte(r.Type()))
	uFlagByte := byte('N')
	if r.UFlag {
		uFlagByte = 'Y'
	}
	buf = append(buf, uFlagByte)
	buf = append(buf, []byte(fmt.Sprintf("%02X", r.HopCount))...)
	buf = append(buf, []byte(fmt.Sprintf("%04X", r.RREQID))...)
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
	buf = append(buf, []byte(fmt.Sprintf("%04X", r.OriginatorSequenceNum))...)
	return buf, nil
}
