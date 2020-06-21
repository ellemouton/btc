package fieldelement

import (
	"errors"
	"fmt"
	"math/big"
)

type FieldElement struct {
	Num   *big.Int
	Prime *big.Int
}

func New(num *big.Int, prime *big.Int) (*FieldElement, error) {
	if num.Cmp(prime) >= 0 || num.Sign() < 0 {
		return nil, errors.New(fmt.Sprintf("Num %v not in field range 0 to %v -1", num, prime))
	}

	return &FieldElement{
		Num:   num,
		Prime: prime,
	}, nil
}

func (f *FieldElement) Copy() *FieldElement {
	num := new(big.Int).Set(f.Num)
	prime := new(big.Int).Set(f.Prime)

	return &FieldElement{
		Num:   num,
		Prime: prime,
	}
}

func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%v(%v)", f.Num, f.Prime)
}

func (f *FieldElement) Hex() string {
	return fmt.Sprintf("%064x", f.Num)
}

func (f *FieldElement) Add(o *FieldElement) (*FieldElement, error) {
	if f.Prime.Cmp(o.Prime) != 0 {
		return nil, errors.New("Cant add numbers of different fields")
	}

	num := &big.Int{}
	num.Add(f.Num, o.Num)
	num.Mod(num, f.Prime)

	if num.Sign() < 0 {
		num.Add(num, f.Prime)
	}

	return &FieldElement{
		Num:   num,
		Prime: f.Prime,
	}, nil
}

func (f *FieldElement) Sub(o *FieldElement) (*FieldElement, error) {
	if f.Prime.Cmp(o.Prime) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	num := &big.Int{}
	num.Sub(f.Num, o.Num)
	num.Mod(num, f.Prime)

	if num.Sign() < 0 {
		num.Add(num, f.Prime)
	}

	return &FieldElement{
		Num:   num,
		Prime: f.Prime,
	}, nil
}

func (f *FieldElement) Mul(o *FieldElement) (*FieldElement, error) {
	if f.Prime.Cmp(o.Prime) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	num := &big.Int{}
	num.Mul(f.Num, o.Num)
	num.Mod(num, f.Prime)

	if num.Sign() < 0 {
		num.Add(num, f.Prime)
	}

	return &FieldElement{
		Num:   num,
		Prime: f.Prime,
	}, nil
}

func (f *FieldElement) Pow(e *big.Int) (*FieldElement, error) {
	p := &big.Int{}
	p.Sub(f.Prime, big.NewInt(1))

	n := &big.Int{}
	n.Mod(e, p)

	if n.Sign() < 0 {
		n.Add(n, p)
	}

	num := &big.Int{}
	num.Exp(f.Num, n, f.Prime)

	if num.Sign() < 0 {
		num.Add(num, f.Prime)
	}

	return &FieldElement{
		Num:   num,
		Prime: f.Prime,
	}, nil
}

func (f *FieldElement) Div(o *FieldElement) (*FieldElement, error) {
	if f.Prime.Cmp(o.Prime) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	exp := &big.Int{}
	exp.Sub(f.Prime, big.NewInt(2))

	n1, err := o.Pow(exp)
	if err != nil {
		return nil, err
	}

	n2, err := f.Mul(n1)
	if err != nil {
		return nil, err
	}

	return n2, nil
}
