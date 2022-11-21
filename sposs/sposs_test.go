package sposs

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/sachaservan/pacl/algebra"
)

// https://datatracker.ietf.org/doc/html/rfc3526#page-3
const primeHexP = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200CBBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
const generatorG = "2"

func TestingGroup() *algebra.Group {
	rand.Seed(time.Now().Unix())

	p := FromSafeHex(primeHexP)
	g := FromSafeHex(generatorG)

	if !p.ProbablyPrime(10) {
		panic("p is not prime")
	}

	// Initialize field values
	baseField := algebra.NewField(p)
	group := algebra.NewGroup(baseField, baseField.NewElement(g))

	return group
}

func TestFullSPoSS(t *testing.T) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	for i := 0; i < 100; i++ {

		x := pp.ExpField.RandomElement()

		// generate additive shares of g^x
		gX := pp.Group.NewElement(x.Int).Value
		additiveShareA, additiveShareB := pp.LinearShares(gX)

		// client proof of knowledge
		proofA, proofB := pp.GenProof(x)

		// step 2: each server uses the received audit share to update the private and public audits
		auditShareA := pp.Audit(additiveShareA, proofA)
		auditShareB := pp.Audit(additiveShareB, proofB)

		// step 3: check that all the values are correct (i.e., the client didn't provide a bad proof)
		okA := pp.CheckAudit(auditShareA, auditShareB)
		okB := pp.CheckAudit(auditShareA, auditShareB)

		if !okA || !okB {
			t.Fatalf("SPoSS audit and verification test failed")
		}
	}
}

func BenchmarkProve(b *testing.B) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	x := pp.Group.Field.RandomElement()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pp.GenProof(x)
	}
}

func BenchmarkAudit(b *testing.B) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	x := pp.Group.Field.RandomElement()
	gX := pp.Group.NewElement(x.Int).Value
	proofA, _ := pp.GenProof(x)
	additiveShareA, _ := pp.LinearShares(gX)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pp.Audit(additiveShareA, proofA)
	}
}

func BenchmarkVerify(b *testing.B) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	x := pp.Group.Field.RandomElement()
	gX := pp.Group.NewElement(x.Int).Value
	proofA, proofB := pp.GenProof(x)
	additiveShareA, additiveShareB := pp.LinearShares(gX)
	auditShareA := pp.Audit(additiveShareA, proofA)
	auditShareB := pp.Audit(additiveShareB, proofB)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pp.CheckAudit(auditShareA, auditShareB)
	}
}

func BenchmarkExp(b *testing.B) {
	group := TestingGroup()

	_, x := group.RandomElement()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		group.NewElement(x)
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
