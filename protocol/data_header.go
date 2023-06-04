package protocol

type CustomHeader struct {
}

func (h *CustomHeader) HeaderSize() int {
	return 1
}
