package messages

import "testing"

func TestRollover(t *testing.T) {
	var incoming, current int16
	incoming, current = 32767, -32768

	if compareSeqnums(current, incoming) {
		t.Fatal("current seqNum should be considered as fresher.")
	}
}

func TestSeqNumDiscard(t *testing.T) {
	var incoming, current int16
	incoming, current = 10, 15

	if compareSeqnums(current, incoming) {
		t.Fatal("The current seqNum should've been bigger.")
	}
}

func TestSeqNumAccept(t *testing.T) {
	var incoming, current int16
	incoming, current = 15, 10

	if !compareSeqnums(current, incoming) {
		t.Fatal("The new seqNum should've been bigger.")
	}
}

func TestSameSeqNum(t *testing.T) {
	var incoming, current int16
	incoming, current = 1, 1
	if !compareSeqnums(current, incoming) {

	}
}
