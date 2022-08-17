package paclsposs

import (
	"bytes"
	"math/big"

	"github.com/sachaservan/pacl/algebra"
	"github.com/sachaservan/pacl/sposs"
	dpf "github.com/sachaservan/vdpf"
)

type ProofShare struct {
	DPFKey      *dpf.DPFKey // DPF or VDPF key
	PrfKey      dpf.PrfKey  // prf used for PRG
	ShareNumber uint
	ProofShare  *sposs.ProofShare // public key (Schnorr) PACL for VDPFs
}

type AuditShare struct {
	Share    *sposs.AuditShare
	Pi       []byte                // VDPF proof
	KeyShare *algebra.FieldElement // for testing purposes
}

func (kl *KeyListParams) NewProof(idx uint64, x *algebra.FieldElement) []*ProofShare {

	if kl.NumKeys == 0 {
		panic("list size is set to zero; something is wrong")
	}

	prfKey := dpf.GeneratePRFKey()

	// initialize the DPF
	pf := dpf.ClientVDPFInitialize(prfKey, [2]dpf.HashKey{kl.HKey1, kl.HKey2})

	// gen the dpf keys
	keyA, keyB := pf.GenVDPFKeys(idx, kl.FSSDomain)

	proofX := new(big.Int).Set(x.Int)

	// flip x based on which server the key is "retrieved" from
	resB := pf.BatchEval(keyB, []uint64{idx})
	if resB[0] == 1 {
		// we need to compute x' such that g^x' = -g^x = p - g^x mod p = -1g^x mod p = g^q+x
		// compute q + x mod 2q
		q := kl.Field.Pminus1()
		q.Div(q, big.NewInt(2))
		proofX = new(big.Int).Add(proofX, q)
		proofX.Mod(proofX, kl.Field.Pminus1())
	}

	spossProofA, spossProofB := kl.ProofPP.GenProof(kl.ProofPP.ExpField.NewElement(proofX))

	shares := make([]*ProofShare, 2)

	// share for server A
	shares[0] = &ProofShare{}
	shares[0].ShareNumber = 0
	shares[0].PrfKey = pf.PrfKey
	shares[0].DPFKey = keyA
	shares[0].ProofShare = spossProofA

	// share for server B
	shares[1] = &ProofShare{}
	shares[1].ShareNumber = 1
	shares[1].PrfKey = pf.PrfKey
	shares[1].DPFKey = keyB
	shares[1].ProofShare = spossProofB

	return shares
}

func (kl *KeyList) Audit(proof *ProofShare) *AuditShare {
	bits, pi := kl.ExpandVDPF(proof)
	return kl.computePrepareAudit(proof, bits, pi)
}

func (kl *KeyList) CheckAudit(auditShares ...*AuditShare) bool {

	vdpfOk := bytes.Equal(auditShares[0].Pi, auditShares[1].Pi)
	spossOk := kl.ProofPP.CheckAudit(auditShares[0].Share, auditShares[1].Share)

	return vdpfOk && spossOk
}

func (kl *KeyList) ExpandVDPF(proof *ProofShare) ([]byte, []byte) {

	var res []byte
	var pi []byte

	pf := dpf.ServerVDPFInitialize(proof.PrfKey, [2]dpf.HashKey{kl.HKey1, kl.HKey2})

	if kl.FullDomain {
		// run the optimized full-domain evaluation strategy
		res, pi = pf.FullDomainVerEval(proof.DPFKey)
	} else {
		res, pi = pf.BatchVerEval(proof.DPFKey, kl.KeyIndices)
	}
	return res, pi
}

// uses the expanded DPF bits to "select" the public key in the keylist
// over which the audit is going to be performed
func (kl *KeyList) computePrepareAudit(proof *ProofShare, bits []byte, pi []byte) *AuditShare {

	// final result
	accumulator := kl.Field.AddIdentity()
	for i := uint64(0); i < kl.NumKeys; i++ {
		if bits[i] == 1 {
			// add result to running sum (mod q)
			kl.Field.AddInplace(accumulator, kl.PublicKeys[i].Value)
		}
	}

	spossAudit := kl.ProofPP.Audit(accumulator, proof.ProofShare)
	return &AuditShare{Share: spossAudit, Pi: pi, KeyShare: accumulator}
}
