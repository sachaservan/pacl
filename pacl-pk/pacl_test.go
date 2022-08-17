package paclpk

import (
	"crypto/elliptic"
	"testing"
)

// test configuration parameters
const TestNumKeys = 512
const TestNumSubkeys = 10 // for inclusion predicate only
const TestFSSDomain = 32
const TestPredicate = Inclusion

const StatSecPar = 128
const NumQueries = 100 // number of queries to run

func TestProveAuditVerify(t *testing.T) {

	for i := 0; i < NumQueries; i++ {
		kl, key, idx := GenerateTestingKeyList(
			TestNumKeys,
			TestFSSDomain,
			elliptic.P256(),
			TestPredicate,
			TestNumSubkeys)

		proofShares := kl.NewProof(idx, key)

		klB := kl.CloneKeyList()
		klB.FlipSignOfKeys()

		auditA := kl.Audit(proofShares[0])
		auditB := klB.Audit(proofShares[1])

		if !kl.CheckAudit(auditA, auditB) {
			t.Fatalf("CheckAudit failed")
		}
	}
}

func BenchmarkBaseline(b *testing.B) {

	numKeys := uint64(1000)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(
		numKeys,
		fssDomain,
		elliptic.P256(),
		TestPredicate,
		TestNumSubkeys)
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kl.ExpandDPF(shares[0])
	}
}

func BenchmarkPACLSingle(b *testing.B) {

	numKeys := uint64(1)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(
		numKeys,
		fssDomain,
		elliptic.P256(),
		TestPredicate,
		TestNumSubkeys)
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		audit := kl.Audit(shares[0])
		kl.CheckAudit(audit, audit)
	}
}

func BenchmarkPACLMany(b *testing.B) {

	numKeys := uint64(1024)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(
		numKeys,
		fssDomain,
		elliptic.P256(),
		TestPredicate,
		TestNumSubkeys)
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		audit := kl.Audit(shares[0])
		kl.CheckAudit(audit, audit)
	}
}
