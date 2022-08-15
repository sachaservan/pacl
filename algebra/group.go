package algebra

import (
	"math/big"
)

type Group struct {
	Field *Field
	G     *FieldElement
}

type GroupElement struct {
	Value *FieldElement
}

// new group over a specified field with generator g
func NewGroup(f *Field, g *FieldElement) *Group {
	return &Group{f, g}
}

// multiply two group elements
func (g *Group) Mul(a, b *GroupElement) *GroupElement {
	newElement := g.Field.Mul(a.Value, b.Value)
	return &GroupElement{newElement}
}

// return multiplicative inverse of the element
func (g *Group) MulInv(a *GroupElement) *GroupElement {
	newElement := g.Field.MulInv(a.Value)
	return &GroupElement{newElement}
}

// new element g**alpha mod P = 2q+1
func (g *Group) NewElement(a *big.Int) *GroupElement {
	newElement := g.Field.Exp(g.G, a)
	return &GroupElement{newElement}
}

// new random element in the group (also returns discrete log)
func (g *Group) RandomElement() (*GroupElement, *big.Int) {
	// should make it not repeat this calculation
	a := randomInt(g.Field.P)

	return g.NewElement(a), a
}

func (g *Group) Identity() *GroupElement {
	return g.NewElement(big.NewInt(0)) // g^0 = 1
}

func (elem *GroupElement) Cmp(b *GroupElement) int {
	return elem.Value.Cmp(b.Value)
}

func (elem *GroupElement) Copy() *GroupElement {
	return &GroupElement{&FieldElement{big.NewInt(0).SetBytes(elem.Value.Int.Bytes())}}
}

func (f *Field) Pminus1() *big.Int {
	return new(big.Int).Sub(f.P, big.NewInt(1))
}
