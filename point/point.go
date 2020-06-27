package point

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ellemouton/btc/fieldelement"
)

type Point interface {
	Add(Point) (Point, error)
	Mul(*big.Int) (Point, error)
	Copy() Point
	String() string
	GetA() fieldelement.FieldElement
	GetB() fieldelement.FieldElement
	GetX() fieldelement.FieldElement
	GetY() fieldelement.FieldElement
}

type point struct {
	X fieldelement.FieldElement
	Y fieldelement.FieldElement
	A fieldelement.FieldElement
	B fieldelement.FieldElement
}

func New(x, y, a, b fieldelement.FieldElement) (Point, error) {
	if x == nil && y == nil {
		return &point{
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

	return &point{
		X: x,
		Y: y,
		A: a,
		B: b,
	}, nil
}

func (p *point) GetA() fieldelement.FieldElement {
	return p.A
}

func (p *point) GetB() fieldelement.FieldElement {
	return p.B
}

func (p *point) GetX() fieldelement.FieldElement {
	return p.X
}

func (p *point) GetY() fieldelement.FieldElement {
	return p.Y
}

func (p *point) Copy() Point {
	a := p.A.Copy()
	b := p.B.Copy()
	x := p.X.Copy()
	y := p.Y.Copy()

	return &point{
		A: a,
		B: b,
		X: x,
		Y: y,
	}
}

func (p *point) String() string {
	return fmt.Sprintf("(%v,%v)_%v_%v", p.X, p.Y, p.A, p.B)
}

func (p *point) Add(o Point) (Point, error) {
	if !reflect.DeepEqual(p.A, o.GetA()) || !reflect.DeepEqual(p.B, o.GetB()) {
		return nil, errors.New(fmt.Sprintf("points %v, %v are not on the same curve", p, o))
	}

	if p.X == nil {
		return &point{
			X: o.GetX(),
			Y: o.GetY(),
			A: o.GetA(),
			B: o.GetB(),
		}, nil
	}

	if o.GetX() == nil {
		return &point{
			X: p.X,
			Y: p.Y,
			A: p.A,
			B: p.B,
		}, nil
	}

	if reflect.DeepEqual(p.X, o.GetX()) && !reflect.DeepEqual(p.Y, o.GetY()) {
		return &point{
			X: nil,
			Y: nil,
			A: p.A,
			B: p.B,
		}, nil
	}

	if reflect.DeepEqual(p.X, o.GetX()) && p.Y.GetNum().Sign() == 0 {
		return &point{
			X: nil,
			Y: nil,
			A: p.A,
			B: p.B,
		}, nil
	}

	if !reflect.DeepEqual(p.X, o.GetX()) {
		s1, err := (o.GetY()).Sub(p.Y)
		if err != nil {
			return nil, err
		}

		s2, err := (o.GetX()).Sub(p.X)
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

		x3, err := t2.Sub(o.GetX())
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

		return &point{
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

	return &point{
		X: x3,
		Y: y3,
		A: p.A,
		B: p.B,
	}, nil
}

func (p *point) Mul(c *big.Int) (Point, error) {
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
