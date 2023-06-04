package protocol

import "encoding/binary"

type RREPHeader struct {
	T          uint8
	Reserved   uint8
	HopCount   uint8
	DestAddr   uint16
	DestSeqNum int16
	OriginAddr uint16
}

func (h *RREPHeader) HeaderSize() int {
	return 8
}

func (h *RREPHeader) Marshal() []byte {
	data := make([]byte, h.HeaderSize())

	// Convert fields to binary representation
	data[0] = (h.T & 0x03) << 6
	data[0] |= h.Reserved
	data[1] = h.HopCount
	binary.BigEndian.PutUint16(data[2:4], h.DestAddr)
	binary.BigEndian.PutUint16(data[4:6], uint16(h.DestSeqNum))
	binary.BigEndian.PutUint16(data[6:8], h.OriginAddr)

	return data
}

func (h *RREPHeader) Unmarshal(data []byte) {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return
	}

	h.T = (data[0] >> 6) & 0x03
	h.Reserved = data[0] & 0x3F
	h.HopCount = data[1]
	h.DestAddr = binary.BigEndian.Uint16(data[2:4])
	h.DestSeqNum = int16(binary.BigEndian.Uint32(data[4:6]))
	h.OriginAddr = binary.BigEndian.Uint16(data[6:8])
}
