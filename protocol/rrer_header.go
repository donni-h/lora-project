package protocol

import "encoding/binary"

type RRERHeader struct {
	T             uint8
	Reserved      uint8
	DestCount     uint8
	UnreachAddr   uint16
	UnreachSeqNum int16
}

func (h *RRERHeader) HeaderSize() int {
	return 6
}

func (h *RRERHeader) Marshal() []byte {
	data := make([]byte, h.HeaderSize())
	data[0] = (h.T & 0x03) << 6
	data[0] |= h.Reserved
	data[1] = h.DestCount
	binary.BigEndian.PutUint16(data[2:4], h.UnreachAddr)
	binary.BigEndian.PutUint16(data[4:6], uint16(h.UnreachSeqNum))

	return data
}

func (h *RRERHeader) Unmarshal(data []byte) {
	if len(data) < h.HeaderSize() {
		return
	}
	h.T = (data[0] >> 6) & 0x03
	h.Reserved = data[0] & 0x3F
	h.DestCount = data[1]
	h.UnreachAddr = binary.BigEndian.Uint16(data[2:4])
	h.UnreachSeqNum = int16(binary.BigEndian.Uint16(data[4:6]))
}
