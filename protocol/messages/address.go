package messages

import (
	"encoding/hex"
	"fmt"
)

type Address [2]byte

func (a Address) String() string {
	return hex.EncodeToString(a[:])
}

func (a Address) MarshalText() ([]byte, error) {
	r := make([]byte, hex.EncodedLen(len(a)))
	hex.Encode(r, a[:])
	return r, nil
}

func (a *Address) UnmarshalText(input []byte) error {
	if expected := hex.EncodedLen(len(a)); expected != len(input) {
		return fmt.Errorf("wrong len, expected %d; got %d", expected, len(input))
	}
	fmt.Println(input)
	_, err := hex.Decode(a[:], input)
	return err
}
