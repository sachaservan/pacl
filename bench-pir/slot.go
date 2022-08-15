package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

// Slot is a set of bytes which can be xor'ed and comapred
type Slot struct {
	Data []byte
}

// XorSlots compute xor a and b storing result in a
func XorSlots(a, b *Slot) {

	if len(a.Data) < len(b.Data) {
		for j := 0; j < len(a.Data); j++ {
			a.Data[j] ^= b.Data[j]
		}
	} else {
		for j := 0; j < len(b.Data); j++ {
			a.Data[j] ^= b.Data[j]
		}
	}
}

// Compare returns the comparison of the two byte arrays
// 0 if slot == other
// -1 if slot < other
// 1 if slot > other
func (slot *Slot) Compare(other *Slot) int {
	return bytes.Compare(slot.Data, other.Data)
}

// NewSlot returns a slot populated with data
func NewSlot(b []byte) *Slot {
	return &Slot{
		Data: b,
	}
}

// NewEmptySlot returns an all-zero slot
func NewEmptySlot(numBytes int) *Slot {
	return &Slot{
		Data: make([]byte, numBytes),
	}
}

// NewRandomSlot returns a slot filled with random bytes
func NewRandomSlot(numBytes int) *Slot {
	slotData := make([]byte, numBytes)
	_, err := rand.Read(slotData)
	if err != nil {
		panic(fmt.Sprintf("Generating random bytes failed with %v\n", err))
	}

	return &Slot{slotData}
}
