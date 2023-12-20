package paclpk

import (
	"crypto/elliptic"
	"math"
	"math/big"
	"math/rand"

	"github.com/sachaservan/pacl/algebra"
	"github.com/sachaservan/pacl/ec"
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
	Curve         *ec.EC
	PredicateType PredicateType
}

type KeyList struct {
	KeyListParams
	PublicKeys []*ec.Point
}

func (kl *KeyList) CloneKeyList() *KeyList {
	clone := KeyList{}
	clone.Curve = kl.Curve
	clone.PublicKeys = make([]*ec.Point, kl.NumKeys)
	clone.NumKeys = kl.NumKeys
	clone.FSSDomain = kl.FSSDomain
	clone.FullDomain = kl.FullDomain
	clone.KeyIndices = kl.KeyIndices
	clone.PredicateType = kl.PredicateType

	for i := uint64(0); i < kl.NumKeys; i++ {
		clone.PublicKeys[i] = kl.PublicKeys[i].Copy()
	}

	return &clone
}

// same as GenerateRandomKeyList but all keys are the same
// this is useful for testing as generating the full list is time consuming
// returns: a key list, a key, and the index of the associated public key
func GenerateTestingKeyList(
	numKeys uint64,
	fssDomain uint,
	curve elliptic.Curve,
	pred PredicateType,
	numSubkeys uint64) (*KeyList, *algebra.FieldElement, uint64) {

	if pred == Inclusion {
		// increase the domain of the DPF to account for the extra
		// evaluations (the subtree that expands to numSubkeys leaves)
		fssDomain += uint(math.Ceil(math.Log2(float64(numSubkeys))))

		// increase the total number of keys to account for the extra subkeys
		// over which the verifiers must select the correct key
		numKeys *= numSubkeys
	}

	c := &ec.EC{Curve: curve, Field: algebra.NewField(curve.Params().N)}
	kl := KeyList{}
	kl.KeyListParams.Curve = c
	kl.PublicKeys = make([]*ec.Point, numKeys)
	kl.NumKeys = numKeys
	kl.FSSDomain = fssDomain
	kl.PredicateType = pred
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	key, gkey, _ := kl.Curve.NewRandomPoint()
	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64() % (1 << fssDomain)
		kl.PublicKeys[i] = gkey.Copy()
	}

	keyElem := kl.Curve.Field.NewElement(new(big.Int).SetBytes(key))
	idx := rand.Uint64() % numKeys

	return &kl, keyElem, kl.KeyIndices[idx]
}

func GenerateBenchmarkKeyList(
	numKeys uint64,
	fssDomain uint,
	curve elliptic.Curve,
	pred PredicateType,
	numSubkeys uint64) (*KeyList, *algebra.FieldElement, uint64, uint64) {

	if pred == Inclusion {
		// increase the domain of the DPF to account for the extra
		// evaluations (the subtree that expands to numSubkeys leaves)
		fssDomain += uint(math.Ceil(math.Log2(float64(numSubkeys))))

		// increase the total number of keys to account for the extra subkeys
		// over which the verifiers must select the correct key
		numKeys *= numSubkeys
	}

	c := &ec.EC{Curve: curve, Field: algebra.NewField(curve.Params().P)}
	kl := KeyList{}
	kl.KeyListParams.Curve = c
	kl.PublicKeys = make([]*ec.Point, numKeys)
	kl.NumKeys = numKeys
	kl.FSSDomain = fssDomain
	kl.PredicateType = pred
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	key, gkey, _ := kl.Curve.NewRandomPoint()
	kl.PublicKeys[0] = gkey
	for i := uint64(1); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.PublicKeys[i] = kl.Curve.Add(kl.PublicKeys[i-1], gkey)
	}

	keyElem := kl.Curve.Field.NewElement(new(big.Int).SetBytes(key))

	idx := rand.Uint64() % numKeys
	return &kl, keyElem, idx, kl.KeyIndices[idx]
}

// sets g^x to g^-x
func (kl *KeyList) FlipSignOfKeys() {
	for i := range kl.PublicKeys {
		kl.PublicKeys[i] = kl.Curve.Inverse(kl.PublicKeys[i])
	}
}

// computes an additive shares in a field that sum to z
func ComputeMaskingShares(f *algebra.Field, z *algebra.FieldElement) []*algebra.FieldElement {
	s1 := f.RandomElement()
	s2 := f.Sub(z, s1)

	res := make([]*algebra.FieldElement, 2)
	res[0] = s1
	res[1] = s2

	return res
}
