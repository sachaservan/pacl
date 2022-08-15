package paclpk

import (
	"github.com/sachaservan/pacl/algebra"
	"github.com/sachaservan/pacl/ec"
	dpf "github.com/sachaservan/vdpf"
)

type PredicateType int

const (
	Equality  PredicateType = 0
	Inclusion PredicateType = 1
)

type ProofShare struct {
	DPFKey      *dpf.DPFKey // DPF key
	PrfKey      dpf.PrfKey  // PRF key used for the PRG in DPF construction
	ShareNumber uint

	// PACL
	PredicateType PredicateType
	KeyShare      *algebra.FieldElement
}

type AuditShare struct {
	Share *ec.Point
}

func (kl *KeyListParams) NewProof(idx uint64, x *algebra.FieldElement) []*ProofShare {

	if kl.NumKeys == 0 {
		panic("list size is set to zero; something is wrong")
	}

	// initialize the DPF
	prfKey := dpf.GeneratePRFKey()
	pf := dpf.ClientDPFInitialize(prfKey)

	// gen the dpf keys
	keyA, keyB := pf.GenDPFKeys(idx, kl.FSSDomain)

	// flip x based on which server the key is "retrieved" from
	resB := pf.BatchEval(keyB, []uint64{idx})
	if resB[0] == 0 {
		x = kl.Curve.Field.Negate(x)
	}

	// secret share the access key x
	keyShares := ComputeMaskingShares(kl.Curve.Field, x)

	// shares provided to each verifier
	shares := make([]*ProofShare, 2)

	// share for verifier A
	shares[0] = &ProofShare{}
	shares[0].ShareNumber = 0
	shares[0].PrfKey = pf.PrfKey
	shares[0].DPFKey = keyA
	shares[0].KeyShare = keyShares[0]

	// share for verifier B
	shares[1] = &ProofShare{}
	shares[1].ShareNumber = 1
	shares[1].PrfKey = pf.PrfKey
	shares[1].DPFKey = keyB
	shares[1].KeyShare = keyShares[1]

	return shares
}

func (kl *KeyList) Audit(proof *ProofShare) *AuditShare {
	bits := kl.ExpandDPF(proof)
	return kl.computeAudit(proof, bits)
}

func (kl *KeyList) CheckAudit(auditShares ...*AuditShare) bool {
	accumulator := auditShares[0].Share
	for i := 1; i < len(auditShares); i++ {
		accumulator = kl.Curve.Add(accumulator, auditShares[i].Share)
	}

	return kl.Curve.IsIdentity(accumulator)
}

func (kl *KeyList) ExpandDPF(proof *ProofShare) []byte {

	pf := dpf.ServerDPFInitialize(proof.PrfKey)

	if kl.FullDomain {
		// run the optimized full-domain evaluation strategy
		return pf.FullDomainEval(proof.DPFKey)
	} else {
		return pf.BatchEval(proof.DPFKey, kl.KeyIndices)
	}
}

// uses the expanded DPF bits to "select" the public key in the keylist
// over which the audit is going to be performed
func (kl *KeyList) computeAudit(proof *ProofShare, bits []byte) *AuditShare {

	// final result
	accumulator, _ := kl.Curve.IdentityPoint()

	for i := uint64(0); i < kl.NumKeys; i++ {
		if bits[i] == 1 {
			// add result to running sum (mod q)
			accumulator = kl.Curve.Add(accumulator, kl.PublicKeys[i])
		}
	}

	share, _ := kl.Curve.NewPoint(proof.KeyShare.Int)
	accumulator = kl.Curve.Add(accumulator, share)

	return &AuditShare{Share: accumulator}
}
