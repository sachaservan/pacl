package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"
)

type Point struct {
	X, Y *big.Int
}

// This is just a bitmask with the number of ones starting at 8 then
// incrementing by index. To account for fields with bitsizes that are not a whole
// number of bytes, we mask off the unnecessary bits. h/t agl
var mask = []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}

// NewRandomPoint: Generates a new random point on the curve specified in curveParams.
func NewRandomPoint(curve elliptic.Curve) ([]byte, *Point, error) {

	s, _, err := RandomCurveScalar(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	x, y := curve.ScalarBaseMult(s)
	return s, &Point{x, y}, nil
}

func RandomCurveScalar(curve elliptic.Curve, rand io.Reader) ([]byte, *big.Int, error) {
	N := curve.Params().N // base point subgroup order
	bitLen := N.BitLen()
	byteLen := (bitLen + 7) >> 3
	buf := make([]byte, byteLen)

	// When in doubt, do what agl does in elliptic.go. Presumably
	// new(big.Int).SetBytes(b).Mod(N) would introduce bias, so we're sampling.
	for {
		_, err := io.ReadFull(rand, buf)
		if err != nil {
			return nil, nil, err
		}
		// Mask to account for field sizes that are not a whole number of bytes.
		buf[0] &= mask[bitLen%8]
		// Check if scalar is in the correct range.
		if new(big.Int).SetBytes(buf).Cmp(N) >= 0 {
			continue
		}
		break
	}

	return buf, new(big.Int).SetBytes(buf), nil
}

func curvePointScalarMult(curve elliptic.Curve, scalar *big.Int) *Point {
	x, y := curve.ScalarBaseMult(scalar.Bytes())
	return &Point{x, y}
}

func curvePointAdd(curve elliptic.Curve, pointA, pointB *Point) *Point {
	x, y := curve.Add(pointA.X, pointA.Y, pointB.X, pointB.Y)
	return &Point{x, y}
}
