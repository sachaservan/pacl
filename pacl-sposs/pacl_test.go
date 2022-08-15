package paclsposs

import (
	"fmt"
	"testing"
)

// test configuration parameters
const TestNumKeys = 512
const TestFSSDomain = 9
const BenchmarkNumKeys = 2000000
const StatSecPar = 128

func TestProveAuditVerify(t *testing.T) {

	group := DefaultGroup()

	kl, key, idx := GenerateTestingKeyList(TestNumKeys, TestFSSDomain, group)

	for i := 0; i < 10; i++ {
		proofShares := kl.NewProof(idx, key)

		klB := kl.CloneKeyList()
		klB.FlipSignOfKeys()

		auditA := kl.Audit(proofShares[0])
		auditB := klB.Audit(proofShares[1])

		resExpected := kl.PublicKeys[idx].Value
		resExpectedAlt := klB.PublicKeys[idx].Value

		recoveredKey := kl.Field.Add(auditA.KeyShare, auditB.KeyShare)
		isValidKeyDPF := recoveredKey.Cmp(resExpected) == 0
		isValidKeyDPFAlt := recoveredKey.Cmp(resExpectedAlt) == 0

		fmt.Printf("isValidKeyDPF = %v, isValidKeyDPFAlt = %v\n", isValidKeyDPF, isValidKeyDPFAlt)

		if !kl.CheckAudit(auditA, auditB) {
			t.Fatalf("CheckAudit failed")
		}
	}
}

func BenchmarkBaseline(b *testing.B) {
	numKeys := uint64(1000)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(numKeys, fssDomain, DefaultGroup())
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kl.ExpandVDPF(shares[0])
	}
}

func BenchmarkPACLSingle(b *testing.B) {

	numKeys := uint64(1)
	fssDomain := uint(32)
	kl, x, _ := GenerateBenchmarkKeyList(numKeys, fssDomain, DefaultGroup())
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
	kl, x, _ := GenerateBenchmarkKeyList(numKeys, fssDomain, DefaultGroup())
	shares := kl.NewProof(0, x)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		audit := kl.Audit(shares[0])
		kl.CheckAudit(audit, audit)
	}
}
