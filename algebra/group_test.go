package algebra

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func setupGroup(p *big.Int, n int) (*Group, []*GroupElement) {
	rand.Seed(time.Now().UnixNano())

	field := NewField(p)
	g := findRandomGenerator(field)
	group := NewGroup(field, g)
	elements := make([]*GroupElement, n)
	for i := 0; i < n; i++ {
		val := big.NewInt(rand.Int63())
		sign := rand.Int() % 2
		if sign == 0 {
			val.Sub(big.NewInt(0), val)
		}
		elements[i] = group.NewElement(val)
	}
	return group, elements
}
func TestAddGroup(t *testing.T) {

	n := 100
	p := big.NewInt(1523) // 1523 is a safe prime
	group, elements := setupGroup(p, n)

	sumInt := group.Field.MulIdentity()
	sum := group.Identity()
	for i := 0; i < n; i++ {
		sum = group.Mul(sum, elements[i])
		sumInt = group.Field.Mul(sumInt, elements[i].Value)
	}

	if sum.Value.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Group element should never be zero!")
	}

	if sum.Value.Int.Cmp(sumInt.Int) != 0 {
		t.Fatalf("Sum over group is not correct. expected: %v, got: %v", sumInt, sum.Value.Int)
	}
}

func TestInverseGroup(t *testing.T) {

	n := 100
	p := big.NewInt(1523) // 1523 is a safe prime
	group, elements := setupGroup(p, n)

	sum := group.Identity()
	sumInv := group.Identity()
	for i := 0; i < n; i++ {
		sum = group.Mul(sum, elements[i])
		sumInv = group.Mul(sumInv, group.MulInv(elements[i]))
	}

	test := group.Mul(sum, sumInv)

	if test.Value.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Group element should never be zero!")
	}

	if test.Cmp(group.Identity()) != 0 {
		t.Fatalf("Multiplicative inverse is incorrect!")
	}
}

// Get all prime factors of a given number n
// taken from https://siongui.github.io/2017/05/09/go-find-all-prime-factors-of-integer-number/
func PrimeFactors(n *big.Int) []*big.Int {

	pfs := make([]*big.Int, 0)
	two := big.NewInt(2)
	zero := big.NewInt(0)

	// Get the number of 2s that divide n
	for new(big.Int).Mod(n, two).Cmp(zero) == 0 {
		pfs = append(pfs, two)
		n.Div(n, two)
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	i := big.NewInt(3)
	for {
		// while i divides n, append i and divide n
		for big.NewInt(0).Mod(n, i).Cmp(zero) == 0 {
			pfs = append(pfs, big.NewInt(0).SetBytes(i.Bytes()))
			n = n.Div(n, i)
		}

		test := big.NewInt(0)
		if test.Mul(i, i).Cmp(n) <= 0 {
			i.Add(i, two)
		} else {
			break
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n.Cmp(two) > 0 {
		pfs = append(pfs, n)
	}

	return pfs
}

func findRandomGenerator(field *Field) *FieldElement {

	found := false
	one := big.NewInt(1)
	factors := PrimeFactors(field.Pminus1())
	g := randomInt(field.P)
	for {

		// test if g is a generator
		for i := 0; i < len(factors); i++ {
			pow := new(big.Int).Div(field.Pminus1(), factors[i])
			if new(big.Int).Exp(g, pow, field.P).Cmp(one) == 0 {
				break
			}
			if i+1 == len(factors) {
				found = true
			}
		}

		if found {
			break
		}

		// try a new candidate
		g = randomInt(field.P)
	}

	return field.NewElement(g)
}
