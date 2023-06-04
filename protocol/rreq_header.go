package protocol

import "encoding/binary"

type RREQHeader struct {
	T            uint8
	Reserved     uint8
	ID           uint16
	HopCount     uint8
	DestAddr     uint16
	DestSeqNum   int16
	OriginAddr   uint16
	OriginSeqNum int16
}

func (h *RREQHeader) HeaderSize() int {
	return 12
}

func (h *RREQHeader) Marshal() []byte {
	data := make([]byte, h.HeaderSize())

	//convert fields to binary representation
	data[0] = (h.T & 0x03) << 6
	data[0] |= h.Reserved
	binary.BigEndian.PutUint16(data[1:3], h.ID)
	data[3] = h.HopCount
	binary.BigEndian.PutUint16(data[4:6], h.DestAddr)
	binary.BigEndian.PutUint16(data[6:8], uint16(h.DestSeqNum))
	binary.BigEndian.PutUint16(data[8:10], h.OriginAddr)
	binary.BigEndian.PutUint16(data[10:12], uint16(h.OriginSeqNum))

	return data
}

func (h *RREQHeader) Unmarshal(data []byte) {
	if len(data) < h.HeaderSize() {
		// Handle error: insufficient data
		return
	}

	// Convert binary representation to fields
	h.T = (data[0] >> 6) & 0x03
	h.Reserved = data[0] & 0x3F
	h.ID = binary.BigEndian.Uint16(data[1:3])
	h.HopCount = data[3]
	h.DestAddr = binary.BigEndian.Uint16(data[4:6])
	h.DestSeqNum = int16(binary.BigEndian.Uint16(data[6:8]))
	h.OriginAddr = binary.BigEndian.Uint16(data[8:10])
	h.OriginSeqNum = int16(binary.BigEndian.Uint16(data[10:12]))
}
