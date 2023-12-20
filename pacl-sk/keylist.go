package paclsk

import (
	"math"
	"math/rand"
)

type PredicateType int

const (
	Equality  PredicateType = 0
	Inclusion PredicateType = 1
)

type KeyListParams struct {
	FullDomain    bool
	NumKeys       uint64
	FSSDomain     uint
	KeyIndices    []uint64
	PredicateType PredicateType
}

type KeyList struct {
	KeyListParams
	Keys         []*Slot
	StatSecurity int // key statistical security (e.g., 128)
}

func GenerateTestingKeyList(
	numKeys uint64,
	fssDomain uint,
	pred PredicateType,
	numSubkeys uint64,
) (*KeyList, *Slot, uint64, uint64) {

	if pred == Inclusion {
		// increase the domain of the DPF to account for the extra
		// evaluations (the subtree that expands to numSubkeys leaves)
		fssDomain += uint(math.Ceil(math.Log2(float64(numSubkeys))))

		// increase the total number of keys to account for the extra subkeys
		// over which the verifiers must select the correct key
		numKeys *= numSubkeys
	}

	kl := KeyList{}
	kl.Keys = make([]*Slot, numKeys)
	kl.NumKeys = numKeys
	kl.StatSecurity = 128
	kl.FSSDomain = fssDomain
	kl.PredicateType = pred
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	slot := NewRandomSlot(kl.StatSecurity / 8)

	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.Keys[i] = NewSlot(slot.Data)
	}

	idx := rand.Uint64() % numKeys

	return &kl, kl.Keys[idx], idx, kl.KeyIndices[idx]
}

func GenerateBenchmarkKeyList(
	numKeys uint64,
	fssDomain uint,
	pred PredicateType,
	numSubkeys uint64,
) (*KeyList, *Slot, uint64) {

	if pred == Inclusion {
		// increase the domain of the DPF to account for the extra
		// evaluations (the subtree that expands to numSubkeys leaves)
		fssDomain += uint(math.Ceil(math.Log2(float64(numSubkeys))))

		// increase the total number of keys to account for the extra subkeys
		// over which the verifiers must select the correct key
		numKeys *= numSubkeys
	}

	kl := KeyList{}
	kl.Keys = make([]*Slot, numKeys)
	kl.NumKeys = numKeys
	kl.StatSecurity = 128
	kl.FSSDomain = fssDomain
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64() % (1 << fssDomain)
		slot := NewRandomSlot(kl.StatSecurity / 8)
		kl.Keys[i] = NewSlot(slot.Data)
	}

	idx := rand.Uint64() % numKeys

	return &kl, kl.Keys[idx], kl.KeyIndices[idx]
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
