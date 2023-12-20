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
const primeHexP = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200CBBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"

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

// same as GenerateTestingKeyList but all keys are the same
// this is useful for testing as generating the full list is time consuming
// returns: a key list, a key, and the index of the associated public key
func GenerateTestingKeyList(
	numKeys uint64,
	fssDomain uint,
	group *algebra.Group,
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
		kl.KeyIndices[i] = rand.Uint64() % (1 << fssDomain)
		kl.PublicKeys[i] = gkey.Copy()
	}

	idx := rand.Uint64() % numKeys

	return &kl, key, idx, kl.KeyIndices[idx]
}

func GenerateBenchmarkKeyList(
	numKeys uint64,
	fssDomain uint,
	group *algebra.Group,
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
		kl.KeyIndices[i] = rand.Uint64() % (1 << fssDomain)
		kl.PublicKeys[i] = group.Mul(kl.PublicKeys[i-1], gkey)
	}

	idx := rand.Uint64() % numKeys

	return &kl, key, idx, kl.KeyIndices[idx]
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
