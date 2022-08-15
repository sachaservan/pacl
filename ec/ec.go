package ec

// partly taken from: https://github.com/privacypass/challenge-bypass-server/blob/master/crypto/curve.go

import (
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"

	"github.com/sachaservan/pacl/algebra"
)

type EC struct {
	Curve elliptic.Curve
	Field *algebra.Field
}

type Point struct {
	X, Y *big.Int
}

func (point *Point) Copy() *Point {
	return &Point{
		X: new(big.Int).SetBytes(point.X.Bytes()),
		Y: new(big.Int).SetBytes(point.Y.Bytes())}
}

// This is just a bitmask with the number of ones starting at 8 then
// incrementing by index. To account for fields with bitsizes that are not a whole
// number of bytes, we mask off the unnecessary bits. h/t agl
var mask = []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}

// GeneratorPoint: the base point of the elliptic curve
func (ec *EC) GeneratorPoint() (*Point, error) {

	s := big.NewInt(1)
	x, y := ec.Curve.ScalarBaseMult(s.Bytes())
	return &Point{x, y}, nil
}

// IdentityPoint: the point at infinity
func (ec *EC) IdentityPoint() (*Point, error) {

	s := big.NewInt(0)
	x, y := ec.Curve.ScalarBaseMult(s.Bytes())
	return &Point{x, y}, nil
}

// NewPoint: Generates a new point on the curve specified in curveParams.
func (ec *EC) NewPoint(s *big.Int) (*Point, error) {
	x, y := ec.Curve.ScalarBaseMult(s.Bytes())
	return &Point{x, y}, nil
}

// NewRandomPoint: Generates a new random point on the curve specified in curveParams.
func (ec *EC) NewRandomPoint() ([]byte, *Point, error) {

	s, _, err := ec.RandomCurveScalar(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	x, y := ec.Curve.ScalarBaseMult(s)
	return s, &Point{x, y}, nil
}

func (ec *EC) RandomCurveScalar(rand io.Reader) ([]byte, *big.Int, error) {
	N := ec.Curve.Params().N // base point subgroup order
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

func (ec *EC) ScalarMult(scalar *big.Int) *Point {
	x, y := ec.Curve.ScalarBaseMult(scalar.Bytes())
	return &Point{X: x, Y: y}
}

func (ec *EC) Add(pointA, pointB *Point) *Point {
	x, y := ec.Curve.Add(pointA.X, pointA.Y, pointB.X, pointB.Y)
	return &Point{X: x, Y: y}
}

func (ec *EC) Inverse(pointA *Point) *Point {
	newPoint := &Point{
		X: new(big.Int).SetBytes(pointA.X.Bytes()),
		Y: new(big.Int).Sub(ec.Curve.Params().P, pointA.Y)}
	return newPoint
}

func (ec *EC) IsEqual(pointA, pointB *Point) bool {
	return pointA.X.Cmp(pointB.X) == 0 && pointA.Y.Cmp(pointB.Y) == 0
}

func (ec *EC) IsIdentity(pointA *Point) bool {
	pointB, _ := ec.IdentityPoint()
	return ec.IsEqual(pointA, pointB)
}
