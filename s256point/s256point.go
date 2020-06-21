package s256point

import (
	"math/big"

	"github.com/ellemouton/btc/point"
	"github.com/ellemouton/btc/s256field"
)

var (
	N *big.Int
	G *S256Point
)

func init() {
	nVal, ok := new(big.Int).SetString(n, 16)
	if !ok {
		panic("invalid hex: " + n)
	}
	N = nVal

	gxVal, ok := new(big.Int).SetString(gx, 16)
	if !ok {
		panic("invalid hex: " + gx)
	}

	gyVal, ok := new(big.Int).SetString(gy, 16)
	if !ok {
		panic("invalid hex: " + gy)
	}

	gPoint, err := New(gxVal, gyVal)
	if err != nil {
		panic("error initializing point G")
	}

	G = gPoint
}

const (
	n  = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	gx = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	gy = "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
)

type S256Point struct {
	*point.Point
}

func New(x, y *big.Int) (*S256Point, error) {
	a, err := s256field.New(big.NewInt(0))
	if err != nil {
		return nil, err
	}

	b, err := s256field.New(big.NewInt(7))
	if err != nil {
		return nil, err
	}

	if x == nil {
		p, err := point.New(nil, nil, a, b)
		return &S256Point{p}, err
	}

	xf, err := s256field.New(x)
	if err != nil {
		return nil, err
	}

	yf, err := s256field.New(y)
	if err != nil {
		return nil, err
	}

	p, err := point.New(xf, yf, a, b)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}

func (s *S256Point) Add(o *S256Point) (*S256Point, error) {
	p, err := s.Point.Add(o.Point)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}

func (s *S256Point) Mul(c *big.Int) (*S256Point, error) {
	coef := &big.Int{}
	coef.Mod(c, N)

	p, err := s.Point.Mul(coef)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}
