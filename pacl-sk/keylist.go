package paclsk

import (
	"math/rand"
)

type KeyListParams struct {
	FullDomain bool
	NumKeys    uint64
	FSSDomain  uint
	KeyIndices []uint64
}

type KeyList struct {
	KeyListParams
	Keys         []*Slot
	StatSecurity int // key statistical security (e.g., 128)
}

func GenerateTestingKeyList(numKeys uint64, fssDomain uint) (*KeyList, *Slot, uint64) {

	kl := KeyList{}
	kl.Keys = make([]*Slot, numKeys)
	kl.NumKeys = numKeys
	kl.StatSecurity = 128
	kl.FSSDomain = fssDomain
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	slot := NewRandomSlot(kl.StatSecurity / 8)

	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.Keys[i] = NewSlot(slot.Data)
	}

	idx := rand.Uint64()

	return &kl, slot, idx
}

func GenerateBenchmarkKeyList(numKeys uint64, fssDomain uint) (*KeyList, *Slot, uint64) {

	kl := KeyList{}
	kl.Keys = make([]*Slot, numKeys)
	kl.NumKeys = numKeys
	kl.StatSecurity = 128
	kl.FSSDomain = fssDomain
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		slot := NewRandomSlot(kl.StatSecurity / 8)
		kl.Keys[i] = NewSlot(slot.Data)
	}

	return &kl, kl.Keys[0], kl.KeyIndices[0]
}

// computes an additive shares in a field that sum to z
func ComputeMaskingShares(z *Slot) []*Slot {
	s1 := NewRandomSlot(len(z.Data))
	s2 := NewEmptySlot(len(z.Data))
	XorSlots(s2, s1)
	XorSlots(s2, z)

	res := make([]*Slot, 2)
	res[0] = s1
	res[1] = s2

	return res
}
