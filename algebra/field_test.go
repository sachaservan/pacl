package algebra

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func setupField(p *big.Int, n int) (*Field, []*FieldElement) {
	rand.Seed(time.Now().UnixNano())

	field := NewField(p)
	elements := make([]*FieldElement, n)
	for i := 0; i < n; i++ {
		val := big.NewInt(rand.Int63())
		sign := rand.Int() % 2
		if sign == 0 {
			val.Sub(big.NewInt(0), val)
		}
		elements[i] = field.NewElement(val)
	}
	return field, elements
}
func TestAddField(t *testing.T) {

	n := 100
	p := big.NewInt(1009)
	field, elements := setupField(p, n)

	sumInt := big.NewInt(0)
	sum := field.AddIdentity()
	for i := 0; i < n; i++ {
		sum = field.Add(sum, elements[i])
		sumInt.Add(sumInt, elements[i].Int)
	}

	sumInt.Mod(sumInt, p)

	if sum.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Sum over field is zero!")
	}

	if sum.Int.Cmp(sumInt) != 0 {
		t.Fatalf("Sum over field is not correct. expected: %v, got: %v", sumInt, sum.Int)
	}
}

func TestSubField(t *testing.T) {

	n := 100
	p := big.NewInt(1009)
	field, elements := setupField(p, n)

	sumInt := big.NewInt(0)
	sum := field.AddIdentity()
	for i := 0; i < n; i++ {
		sum = field.Sub(sum, elements[i])
		sumInt.Sub(sumInt, elements[i].Int)
	}

	sumInt.Mod(sumInt, p)

	if sum.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Subtraction over field is zero!")
	}

	if sum.Int.Cmp(sumInt) != 0 {
		t.Fatalf("Subtraction over field is not correct. expected: %v, got: %v", sum.Int, sumInt)
	}
}

func TestMulField(t *testing.T) {

	n := 500
	p := big.NewInt(1009)
	field, elements := setupField(p, n)

	prodInt := big.NewInt(1)
	prod := field.MulIdentity()
	for i := 0; i < n; i++ {
		if !field.IsZero(elements[i]) {
			prod = field.Mul(prod, elements[i])
			prodInt.Mul(prodInt, elements[i].Int)
		}

		// in case mod P became zero
		if field.IsZero(prod) {
			prodInt = big.NewInt(1)
			prod = field.MulIdentity()
		}
	}

	prodInt.Mod(prodInt, p)

	if prod.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Prod over field is zero")
	}

	if prod.Int.Cmp(prodInt) != 0 {
		t.Fatalf("Prod over field is not correct. expected: %v, got: %v", prod.Int, prodInt)
	}
}

func TestMulInverseField(t *testing.T) {

	n := 500
	p := big.NewInt(1009)
	field, elements := setupField(p, n)

	prod := field.MulIdentity()
	prodInv := field.MulIdentity()
	for i := 0; i < n; i++ {
		prod = field.Mul(prod, elements[i])
		prodInv = field.Mul(prodInv, field.MulInv(elements[i]))

		// in case mod P became zero
		if field.IsZero(prod) || field.IsZero(prodInv) {
			prod = field.MulIdentity()
			prodInv = field.MulIdentity()
		}
	}

	res := field.Mul(prod, prodInv)

	if res.Cmp(field.MulIdentity()) != 0 {
		t.Fatalf("x * (x^-1). Expected: 1, got: %v", res.Int)
	}
}
