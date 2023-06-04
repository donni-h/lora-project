package protocol

type Header interface {
	HeaderSize() int
	Marshal() []byte
	Unmarshal(data []byte)
}
