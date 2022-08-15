package sposs

import (
	"bytes"
	"crypto/sha256"
	"math/big"

	"github.com/sachaservan/pacl/algebra"
)

type PublicParams struct {
	Group    *algebra.Group
	ExpField *algebra.Field
	RandSeed *algebra.FieldElement
}

type ProofShare struct {
	ServerNumber int
	ShareX       *algebra.FieldElement // additive share of x such that g^x =
	ShareU       *algebra.FieldElement // a or b
	ShareC       *algebra.FieldElement // [c]
	D            *algebra.FieldElement // beaver mult opening
	E            *algebra.FieldElement // beaver mult opening
	R            *algebra.FieldElement // randomness
	Nonce        *algebra.FieldElement // nonce used in random oracle
}

type AuditShare struct {
	// in the two verifier case, both verifiers have subtractive shares of zero
	// so we can hash down to save bandwidth
	HashedData [32]byte
}

func NewPublicParams(g *algebra.Group) *PublicParams {
	f := algebra.NewField(g.Field.Pminus1())
	return &PublicParams{g, f, nil}
}

func (pp *PublicParams) GenProof(x *algebra.FieldElement) (*ProofShare, *ProofShare) {

	// generate (additive) secret shares of x
	xA, xB := pp.ExpLinearShares(x)

	// a and b of the beaver triple
	a, b := pp.Group.Field.RandomElement(), pp.Group.Field.RandomElement()

	// shares of c = ab
	ab := pp.Group.Field.Mul(a, b)
	cA, cB := pp.LinearShares(ab)

	// nonces for Fiat-Shamir over secret shares
	nonceA := pp.Group.Field.RandomElement()
	nonceB := pp.Group.Field.RandomElement()

	// compute randomness by applying Fiat-Shamir
	rA := pp.RandomOracle(nonceA, xA, a, cA)
	rB := pp.RandomOracle(nonceB, xB, b, cB)
	r := pp.Group.Field.Add(rA, rB)

	// compute g^[x]
	gxA := pp.Group.NewElement(xA.Int)
	gxB := pp.Group.NewElement(xB.Int)

	// d = rg^xA - a
	d := pp.Group.Field.Mul(r, gxA.Value)
	d = pp.Group.Field.Sub(d, a)

	// e = g^xB - b
	e := pp.Group.Field.Sub(gxB.Value, b)

	return &ProofShare{0, xA, a, cA, d, e, r, nonceA}, &ProofShare{1, xB, b, cB, d, e, r, nonceB}
}

func (pp *PublicParams) Audit(yShare *algebra.FieldElement, proofShare *ProofShare) *AuditShare {

	// recompute the randomness
	r := pp.RandomOracle(proofShare.Nonce, proofShare.ShareX, proofShare.ShareU, proofShare.ShareC)

	// recompute g^[x]
	gx := pp.Group.NewElement(proofShare.ShareX.Int)

	// 2^-1 mod p
	twoInv := pp.Group.Field.MulInv(pp.Group.Field.NewElement(big.NewInt(2)))

	var u *algebra.FieldElement // either d or e depending on the verifier

	var shareV *algebra.FieldElement
	var shareW *algebra.FieldElement
	if proofShare.ServerNumber == 0 {
		u = pp.Group.Field.Mul(proofShare.R, gx.Value)
		u = pp.Group.Field.Sub(u, proofShare.ShareU)
		u = pp.Group.Field.Sub(u, proofShare.D)

		shareV = pp.Group.Field.Mul(proofShare.D, proofShare.E)
		shareV = pp.Group.Field.Mul(shareV, twoInv)
		eu := pp.Group.Field.Mul(proofShare.E, proofShare.ShareU)
		shareV = pp.Group.Field.Add(shareV, eu)
		shareV = pp.Group.Field.Add(shareV, proofShare.ShareC)
		ry := pp.Group.Field.Mul(proofShare.R, yShare)
		shareW = pp.Group.Field.Sub(shareV, ry)

		// compute [r] - r
		r = pp.Group.Field.Sub(r, proofShare.R)

	} else {
		u = pp.Group.Field.Sub(gx.Value, proofShare.ShareU)
		u = pp.Group.Field.Sub(u, proofShare.E)

		shareV = pp.Group.Field.Mul(proofShare.D, proofShare.E)
		shareV = pp.Group.Field.Mul(shareV, twoInv)
		du := pp.Group.Field.Mul(proofShare.D, proofShare.ShareU)
		shareV = pp.Group.Field.Add(shareV, du)
		shareV = pp.Group.Field.Add(shareV, proofShare.ShareC)
		ry := pp.Group.Field.Mul(proofShare.R, yShare)
		shareW = pp.Group.Field.Sub(shareV, ry)

		// turn the additive shares into subtractive shares
		shareW = pp.Group.Field.Negate(shareW)
		u = pp.Group.Field.Negate(u)
		r = pp.Group.Field.Negate(r)
	}

	data := []byte{}
	data = append(data, shareW.Int.Bytes()...)
	data = append(data, r.Int.Bytes()...)
	data = append(data, u.Int.Bytes()...)
	return &AuditShare{sha256.Sum256(data)}
}

func (pp *PublicParams) CheckAudit(auditShareA, auditShareB *AuditShare) bool {
	return bytes.Equal(auditShareA.HashedData[:], auditShareB.HashedData[:])
}

func (pp *PublicParams) RandomOracle(nonceShare, xShare, uShare, cShare *algebra.FieldElement) *algebra.FieldElement {

	data := []byte{}
	data = append(data, nonceShare.Int.Bytes()...)
	data = append(data, xShare.Int.Bytes()...)
	data = append(data, uShare.Int.Bytes()...)
	data = append(data, cShare.Int.Bytes()...)
	bytes := sha256.Sum256(data)
	return pp.Group.Field.NewElement(new(big.Int).SetBytes(bytes[:]))
}

// Return a pair of linear shares for toShare, s.t. share1 + share2 = toShare
func (pp *PublicParams) LinearShares(
	toShare *algebra.FieldElement) (*algebra.FieldElement, *algebra.FieldElement) {
	share1 := pp.Group.Field.RandomElement()
	share2 := pp.Group.Field.Sub(toShare, share1)
	return share1, share2
}

// Return a pair of linear shares for toShare, s.t. share1 + share2 = toShare
// the field is the *exponent field* of the group
func (pp *PublicParams) ExpLinearShares(
	toShare *algebra.FieldElement) (*algebra.FieldElement, *algebra.FieldElement) {
	share1 := pp.ExpField.RandomElement()
	share2 := pp.ExpField.Sub(toShare, share1)
	return share1, share2
}
