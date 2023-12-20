package paclsk

import (
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
		kl, key, _, keyIdx := GenerateTestingKeyList(
			TestNumKeys, TestFSSDomain, TestPredicate, TestNumSubkeys)
		proofShares := kl.NewProof(keyIdx, key)

		auditA := kl.Audit(proofShares[0])
		auditB := kl.Audit(proofShares[1])

		if !kl.CheckAudit(auditA, auditB) {
			t.Fatalf("CheckAudit failed")
		}
	}
}

func BenchmarkBaseline(b *testing.B) {

	numKeys := uint64(1000)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(
		numKeys, fssDomain, TestPredicate, TestNumSubkeys)
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
		numKeys, fssDomain, TestPredicate, TestNumSubkeys)
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		audit := kl.Audit(shares[0])
		kl.CheckAudit(audit, audit)
	}
}

func BenchmarkPACLMany(b *testing.B) {

	numKeys := uint64(1000)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(
		numKeys, fssDomain, TestPredicate, TestNumSubkeys)
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		audit := kl.Audit(shares[0])
		kl.CheckAudit(audit, audit)
	}
}
