package point

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ellemouton/btc/fieldelement"
)

type Point struct {
	X *fieldelement.FieldElement
	Y *fieldelement.FieldElement
	A *fieldelement.FieldElement
	B *fieldelement.FieldElement
}

func New(x, y, a, b *fieldelement.FieldElement) (*Point, error) {
	if x == nil && y == nil {
		return &Point{
			X: nil,
			Y: nil,
			A: a,
			B: b,
		}, nil
	}

	t1, err := y.Pow(big.NewInt(2))
	if err != nil {
		return nil, err
	}

	t2, err := x.Pow(big.NewInt(3))
	if err != nil {
		return nil, err
	}

	t3, err := a.Mul(x)
	if err != nil {
		return nil, err
	}

	t4, err := t2.Add(t3)
	if err != nil {
		return nil, err
	}

	t5, err := t4.Add(b)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(t1, t5) {
		return nil, errors.New(fmt.Sprintf("(%d, %d) is not on the curve", x, y))
	}

	return &Point{
		X: x,
		Y: y,
		A: a,
		B: b,
	}, nil
}

func (p *Point) Copy() *Point {
	a := p.A.Copy()
	b := p.B.Copy()
	x := p.X.Copy()
	y := p.Y.Copy()

	return &Point{
		A: a,
		B: b,
		X: x,
		Y: y,
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)_%v_%v", p.X, p.Y, p.A, p.B)
}

func (p *Point) Add(o *Point) (*Point, error) {
	if !reflect.DeepEqual(p.A, o.A) || !reflect.DeepEqual(p.B, o.B) {
		return nil, errors.New(fmt.Sprintf("points %v, %v are not on the same curve", p, o))
	}

	if p.X == nil {
		return &Point{
			X: o.X,
			Y: o.Y,
			A: o.A,
			B: o.B,
		}, nil
	}

	if o.X == nil {
		return &Point{
			X: p.X,
			Y: p.Y,
			A: p.A,
			B: p.B,
		}, nil
	}

	if reflect.DeepEqual(p.X, o.X) && !reflect.DeepEqual(p.Y, o.Y) {
		return &Point{
			X: nil,
			Y: nil,
			A: p.A,
			B: p.B,
		}, nil
	}

	if reflect.DeepEqual(p.X, o.X) && p.Y.Num.Sign() == 0 {
		return &Point{
			X: nil,
			Y: nil,
			A: p.A,
			B: p.B,
		}, nil
	}

	if !reflect.DeepEqual(p.X, o.X) {
		s1, err := (o.Y).Sub(p.Y)
		if err != nil {
			return nil, err
		}

		s2, err := (o.X).Sub(p.X)
		if err != nil {
			return nil, err
		}

		s, err := s1.Div(s2)
		if err != nil {
			return nil, err
		}

		t1, err := s.Pow(big.NewInt(2))
		if err != nil {
			return nil, err
		}

		t2, err := t1.Sub(p.X)
		if err != nil {
			return nil, err
		}

		x3, err := t2.Sub(o.X)
		if err != nil {
			return nil, err
		}

		t3, err := (p.X).Sub(x3)
		if err != nil {
			return nil, err
		}

		t4, err := s.Mul(t3)
		if err != nil {
			return nil, err
		}

		y3, err := t4.Sub(p.Y)
		if err != nil {
			return nil, err
		}

		return &Point{
			X: x3,
			Y: y3,
			A: p.A,
			B: p.B,
		}, nil
	}

	// p == o
	r, err := (p.X).Pow(big.NewInt(2))
	if err != nil {
		return nil, err
	}

	r1, err := r.Add(r)
	if err != nil {
		return nil, err
	}

	r2, err := r1.Add(r)
	if err != nil {
		return nil, err
	}

	r3, err := r2.Add(p.A)
	if err != nil {
		return nil, err
	}

	n, err := (p.Y).Add(p.Y)
	if err != nil {
		return nil, err
	}

	s, err := r3.Div(n)
	if err != nil {
		return nil, err
	}

	s2, err := s.Pow(big.NewInt(2))
	if err != nil {
		return nil, err
	}

	g, err := (p.X).Add(p.X)
	if err != nil {
		return nil, err
	}

	x3, err := s2.Sub(g)
	if err != nil {
		return nil, err
	}

	d, err := (p.X).Sub(x3)
	if err != nil {
		return nil, err
	}

	d1, err := s.Mul(d)
	if err != nil {
		return nil, err
	}

	y3, err := d1.Sub(p.Y)
	if err != nil {
		return nil, err
	}

	return &Point{
		X: x3,
		Y: y3,
		A: p.A,
		B: p.B,
	}, nil
}

func (p *Point) Mul(c *big.Int) (*Point, error) {
	coef := new(big.Int).Set(c)
	current := p.Copy()

	result, err := New(nil, nil, p.A, p.B)
	if err != nil {
		return nil, err
	}

	for coef.Sign() > 0 {
		z := &big.Int{}
		if z.And(coef, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			result, err = result.Add(current)
			if err != nil {
				return nil, err
			}
		}

		current, err = current.Add(current)
		if err != nil {
			return nil, err
		}

		coef.Rsh(coef, 1)
	}

	return result, nil
}
