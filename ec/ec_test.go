package ec

import (
	"crypto/elliptic"
	"testing"

	"github.com/sachaservan/pacl/algebra"
)

func TestIdentity(t *testing.T) {

	ec := &EC{elliptic.P224(), algebra.NewField(elliptic.P224().Params().P)}
	id, _ := ec.IdentityPoint()
	_, r, _ := ec.NewRandomPoint()
	sum := ec.Add(id, r)

	if !ec.IsEqual(r, sum) {
		t.Fatalf("Identity point is not correct")
	}
}

func TestAdd(t *testing.T) {

	ec := &EC{elliptic.P224(), algebra.NewField(elliptic.P224().Params().P)}
	_, r1, _ := ec.NewRandomPoint()
	_, r2, _ := ec.NewRandomPoint()

	x, y := ec.Curve.Add(r1.X, r1.Y, r2.X, r2.Y)
	p := &Point{
		X: x,
		Y: y,
	}

	sum := ec.Add(r1, r2)

	if !ec.IsEqual(p, sum) {
		t.Fatalf("Add is wrong")
	}
}

func TestInverse(t *testing.T) {

	ec := &EC{elliptic.P224(), algebra.NewField(elliptic.P224().Params().P)}
	_, r, _ := ec.NewRandomPoint()
	id, _ := ec.IdentityPoint()

	inv := ec.Inverse(r)
	res := ec.Add(r, inv)

	if !ec.IsEqual(res, id) {
		t.Fatalf("Inverse is wrong")
	}
}

func BenchmarkCurveAddition(b *testing.B) {

	ec := &EC{elliptic.P224(), algebra.NewField(elliptic.P224().Params().P)}

	list := make([]*Point, 1000)
	for i := 0; i < 1000; i++ {
		_, r, _ := ec.NewRandomPoint()
		list[i] = r
	}

	b.ResetTimer()

	next := 0
	for i := 0; i < b.N; i++ {
		if next+2 > len(list) {
			next = 0
		}

		ec.Add(list[next], list[next+1])
		next++
	}
}
