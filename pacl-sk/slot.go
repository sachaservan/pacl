package paclsk

import (
	"crypto/rand"
	"fmt"
)

// Slot is a set of bytes which can be xor'ed and compared
type Slot struct {
	Data []byte
}

// NewSlotFromString converts a string to a slot type
func NewSlotFromString(s string, slotSize int) *Slot {
	b := []byte(s)
	for i := 0; i < (slotSize - len(s)); i++ {
		b = append(b, 0)
	}
	return &Slot{
		Data: b,
	}
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

// returns true if slot == other
func (slot *Slot) Equal(other *Slot) bool {

	if slot == nil || other == nil {
		return false
	}

	if len(slot.Data) != len(other.Data) {
		return false
	}

	for j := 0; j < len(other.Data); j++ {
		if slot.Data[j] != other.Data[j] {
			return false
		}
	}

	return true
}
