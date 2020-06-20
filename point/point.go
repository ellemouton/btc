package point

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ellemouton/btc/fieldelement"
)

type Point struct {
	x *fieldelement.FieldElement
	y *fieldelement.FieldElement
	a *fieldelement.FieldElement
	b *fieldelement.FieldElement
}

func New(x, y, a, b *fieldelement.FieldElement) (*Point, error) {
	if x == nil && y == nil {
		return &Point{
			x: nil,
			y: nil,
			a: a,
			b: b,
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
		x: x,
		y: y,
		a: a,
		b: b,
	}, nil
}

func (p *Point) Copy() *Point {
	a := p.a.Copy()
	b := p.b.Copy()
	x := p.x.Copy()
	y := p.y.Copy()

	return &Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)_%v_%v", p.x, p.y, p.a, p.b)
}

func (p *Point) Add(o *Point) (*Point, error) {
	if !reflect.DeepEqual(p.a, o.a) || !reflect.DeepEqual(p.b, o.b) {
		return nil, errors.New(fmt.Sprintf("points %v, %v are not on the same curve", p, o))
	}

	if p.x == nil {
		return &Point{
			x: o.x,
			y: o.y,
			a: o.a,
			b: o.b,
		}, nil
	}

	if o.x == nil {
		return &Point{
			x: p.x,
			y: p.y,
			a: p.a,
			b: p.b,
		}, nil
	}

	if reflect.DeepEqual(p.x, o.x) && !reflect.DeepEqual(p.y, o.y) {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}, nil
	}

	if reflect.DeepEqual(p.x, o.x) && p.y.Num.Sign() == 0 {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}, nil
	}

	if !reflect.DeepEqual(p.x, o.x) {
		s1, err := (o.y).Sub(p.y)
		if err != nil {
			return nil, err
		}

		s2, err := (o.x).Sub(p.x)
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

		t2, err := t1.Sub(p.x)
		if err != nil {
			return nil, err
		}

		x3, err := t2.Sub(o.x)
		if err != nil {
			return nil, err
		}

		t3, err := (p.x).Sub(x3)
		if err != nil {
			return nil, err
		}

		t4, err := s.Mul(t3)
		if err != nil {
			return nil, err
		}

		y3, err := t4.Sub(p.y)
		if err != nil {
			return nil, err
		}

		return &Point{
			x: x3,
			y: y3,
			a: p.a,
			b: p.b,
		}, nil
	}

	// p == o
	r, err := (p.x).Pow(big.NewInt(2))
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

	r3, err := r2.Add(p.a)
	if err != nil {
		return nil, err
	}

	n, err := (p.y).Add(p.y)
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

	g, err := (p.x).Add(p.x)
	if err != nil {
		return nil, err
	}

	x3, err := s2.Sub(g)
	if err != nil {
		return nil, err
	}

	d, err := (p.x).Sub(x3)
	if err != nil {
		return nil, err
	}

	d1, err := s.Mul(d)
	if err != nil {
		return nil, err
	}

	y3, err := d1.Sub(p.y)
	if err != nil {
		return nil, err
	}

	return &Point{
		x: x3,
		y: y3,
		a: p.a,
		b: p.b,
	}, nil
}

func (p *Point) Mul(c *big.Int) (*Point, error) {
	coef := new(big.Int).Set(c)
	current := p.Copy()

	result, err := New(nil, nil, p.a, p.b)
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
