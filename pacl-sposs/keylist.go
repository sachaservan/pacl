package paclsposs

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/sachaservan/pacl/algebra"
	"github.com/sachaservan/pacl/sposs"
	dpf "github.com/sachaservan/vdpf"
)

// https://datatracker.ietf.org/doc/html/rfc3526#page-3
const primeHexP = "88426e468d0e90c43ac3d7ff2713ec3e341b1ff2dbdc0f9ef8e7067e5e95d73ab553ffb19d094cae390bb2f1e0c28c4cbbaf3858f071568b120b10a36c9d058b5a219e5842a8ac8c59c8a787b353322e26ee80275fb0d6b39133d7250b9dbd570ea457ad766539196dd93017ecb117e65590422ac309415931554b0e71d6b96008f216782f082cbddfdb7f79b37ace203da13cfe072df9291501efd0edd280c739a7e01010e8782e78ebc556ce7c2a4b54c338d4ee5cc5e2fb668ba6d0a793ea345559768ea104b1b984118b47ea2e8670f722db9d6cdb0e802b79b0c1daa48160308bda2bba41adcc2b884a31a6274be34e11bda421dde626de94a1dc522d47"

// generator g should have a generating subgroup of order 2q
// (NOT the group of quadratic residues as commonly done)
const generatorG = "5"

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
	HKey1         dpf.HashKey    // hash key for VDPF (should be chosen by the verifiers, not the prover)
	HKey2         dpf.HashKey    // hash key for VDPF (should be chosen by the verifiers, not the prover)
	Group         *algebra.Group // multiplicative group of order q
	Field         *algebra.Field // field of order p (elements of Group live in Field)
	ProofPP       *sposs.PublicParams
	PredicateType PredicateType
}

type KeyList struct {
	KeyListParams
	PublicKeys []*algebra.GroupElement
}

func DefaultGroup() *algebra.Group {
	rand.Seed(time.Now().Unix())

	p := FromSafeHex(primeHexP)
	g := FromSafeHex(generatorG)

	// Initialize field values
	baseField := algebra.NewField(p)

	group := algebra.NewGroup(baseField, baseField.NewElement(g))

	if !p.ProbablyPrime(10) {
		panic("p is not prime")
	}

	q := group.Field.Pminus1()
	q.Div(q, big.NewInt(2))

	if big.NewInt(0).Exp(g, q, p).Cmp(big.NewInt(1)) == 0 {
		panic("g isn't a generator of order 2q")
	}

	return group
}

// generate a KeyList of size 'numKeys' where
// each key is a random group element g**(alpha mod q) and where 0 <= alpha <= q-1
func GenerateRandomKeyList(
	numKeys uint64,
	fssDomain uint,
	group *algebra.Group,
	pred PredicateType,
	numSubkeys uint64) *KeyList {

	if pred == Inclusion {
		// increase the domain of the DPF to account for the extra
		// evaluations (the subtree that expands to numSubkeys leaves)
		fssDomain += uint(math.Ceil(math.Log2(float64(numSubkeys))))

		// increase the total number of keys to account for the extra subkeys
		// over which the verifiers must select the correct key
		numKeys *= numSubkeys
	}

	kl := KeyList{}
	kl.PublicKeys = make([]*algebra.GroupElement, numKeys)
	kl.NumKeys = numKeys
	kl.Group = group
	kl.Field = group.Field
	kl.FSSDomain = fssDomain
	kl.PredicateType = pred
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	// for every row, create a random element
	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.PublicKeys[i], _ = group.RandomElement()
	}

	return &kl
}

// same as GenerateRandomKeyList but all keys are the same
// this is useful for testing as generating the full list is time consuming
// returns: a key list, a key, and the index of the associated public key
func GenerateTestingKeyList(
	numKeys uint64,
	fssDomain uint,
	group *algebra.Group,
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

	kl := KeyList{}
	kl.PublicKeys = make([]*algebra.GroupElement, numKeys)
	kl.NumKeys = numKeys
	kl.Group = group
	kl.Field = group.Field
	kl.FSSDomain = fssDomain
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	pp := sposs.NewPublicParams(group)
	kl.ProofPP = pp

	key := kl.ProofPP.ExpField.RandomElement()
	gkey := group.NewElement(key.Int)
	for i := uint64(0); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.PublicKeys[i] = gkey.Copy()
	}

	return &kl, key, 0
}

func GenerateBenchmarkKeyList(
	numKeys uint64,
	fssDomain uint,
	group *algebra.Group,
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

	kl := KeyList{}
	kl.PublicKeys = make([]*algebra.GroupElement, numKeys)
	kl.NumKeys = numKeys
	kl.Group = group
	kl.Field = group.Field
	kl.PredicateType = pred
	kl.FSSDomain = fssDomain
	kl.KeyIndices = make([]uint64, numKeys)
	kl.FullDomain = (1<<fssDomain == numKeys) // only applies when domain = #keys

	pp := sposs.NewPublicParams(group)
	kl.ProofPP = pp

	key := kl.ProofPP.ExpField.RandomElement()
	gkey := group.NewElement(key.Int)
	kl.PublicKeys[0] = gkey
	for i := uint64(1); i < numKeys; i++ {
		kl.KeyIndices[i] = rand.Uint64()
		kl.PublicKeys[i] = group.Mul(kl.PublicKeys[i-1], gkey)
	}

	return &kl, key, 0
}

func (kl *KeyList) CloneKeyList() *KeyList {
	clone := KeyList{}
	clone.PublicKeys = make([]*algebra.GroupElement, kl.NumKeys)
	clone.ProofPP = kl.ProofPP
	clone.NumKeys = kl.NumKeys
	clone.Group = kl.Group
	clone.Field = kl.Group.Field
	clone.HKey1 = kl.HKey1
	clone.HKey2 = kl.HKey2
	clone.FSSDomain = kl.FSSDomain
	clone.FullDomain = kl.FullDomain
	clone.KeyIndices = kl.KeyIndices
	clone.PredicateType = kl.PredicateType

	for i := uint64(0); i < kl.NumKeys; i++ {
		clone.PublicKeys[i] = kl.PublicKeys[i].Copy()
	}

	return &clone
}

// sets g^x to -g^x = p-g^x
func (kl *KeyList) FlipSignOfKeys() {
	for i, k := range kl.PublicKeys {
		newVal := kl.Field.Sub(kl.Field.NewElement(kl.Field.P), k.Value)
		kl.PublicKeys[i] = &algebra.GroupElement{Value: newVal}
	}
}

// from https://github.com/didiercrunch/elgamal/blob/master/elgamal_test.go
func fromHex(hex string) (*big.Int, error) {
	n, err := new(big.Int).SetString(hex, 16)
	if !err {
		msg := fmt.Sprintf("Cannot convert %s to int as hexadecimal", hex)
		return nil, errors.New(msg)
	}
	return n, nil
}

// from https://github.com/didiercrunch/elgamal/blob/master/elgamal_test.go
func FromSafeHex(s string) *big.Int {
	ret, err := fromHex(s)
	if err != nil {
		panic(err)
	}
	return ret
}
