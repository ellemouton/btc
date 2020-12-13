package fieldelement

import (
	"errors"
	"fmt"
	"math/big"
)

type FieldElement interface {
	GetNum() *big.Int
	GetPrime() *big.Int
	Copy() FieldElement
	String() string
	Hex() string
	Add(FieldElement) (FieldElement, error)
	Sub(FieldElement) (FieldElement, error)
	Mul(FieldElement) (FieldElement, error)
	Div(FieldElement) (FieldElement, error)
	Pow(*big.Int) (FieldElement, error)
}

type Element struct {
	Num   *big.Int
	Prime *big.Int
}

func New(num *big.Int, prime *big.Int) (FieldElement, error) {
	if num.Cmp(prime) >= 0 || num.Sign() < 0 {
		return nil, errors.New(fmt.Sprintf("Num %v not in field range 0 to %v -1", num, prime))
	}

	return &Element{
		Num:   num,
		Prime: prime,
	}, nil
}

func (e *Element) GetNum() *big.Int {
	return e.Num
}

func (e *Element) GetPrime() *big.Int {
	return e.Prime
}

func (e *Element) Copy() FieldElement {
	num := new(big.Int).Set(e.Num)
	prime := new(big.Int).Set(e.Prime)

	return &Element{
		Num:   num,
		Prime: prime,
	}
}

func (e *Element) String() string {
	return fmt.Sprintf("FieldElement_%v(%v)", e.Prime, e.Num)
}

func (e *Element) Hex() string {
	return fmt.Sprintf("%064x", e.Num)
}

func (e *Element) Add(o FieldElement) (FieldElement, error) {
	if e.Prime.Cmp(o.GetPrime()) != 0 {
		return nil, errors.New("Cant add numbers of different fields")
	}

	num := &big.Int{}
	num.Add(e.Num, o.GetNum())
	num.Mod(num, e.Prime)

	if num.Sign() < 0 {
		num.Add(num, e.Prime)
	}

	return &Element{
		Num:   num,
		Prime: e.Prime,
	}, nil
}

func (e *Element) Sub(o FieldElement) (FieldElement, error) {
	if e.Prime.Cmp(o.GetPrime()) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	num := &big.Int{}
	num.Sub(e.Num, o.GetNum())
	num.Mod(num, e.Prime)

	if num.Sign() < 0 {
		num.Add(num, e.Prime)
	}

	return &Element{
		Num:   num,
		Prime: e.Prime,
	}, nil
}

func (e *Element) Mul(o FieldElement) (FieldElement, error) {
	if e.Prime.Cmp(o.GetPrime()) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	num := &big.Int{}
	num.Mul(e.Num, o.GetNum())
	num.Mod(num, e.Prime)

	if num.Sign() < 0 {
		num.Add(num, e.Prime)
	}

	return &Element{
		Num:   num,
		Prime: e.Prime,
	}, nil
}

func (e *Element) Pow(exp *big.Int) (FieldElement, error) {
	p := &big.Int{}
	p.Sub(e.Prime, big.NewInt(1))

	n := &big.Int{}
	n.Mod(exp, p)

	if n.Sign() < 0 {
		n.Add(n, p)
	}

	num := &big.Int{}
	num.Exp(e.Num, n, e.Prime)

	if num.Sign() < 0 {
		num.Add(num, e.Prime)
	}

	return &Element{
		Num:   num,
		Prime: e.Prime,
	}, nil
}

func (e *Element) Div(o FieldElement) (FieldElement, error) {
	if e.Prime.Cmp(o.GetPrime()) != 0 {
		return nil, errors.New("Cant subtract numbers of different fields")
	}

	exp := &big.Int{}
	exp.Sub(e.Prime, big.NewInt(2))

	n1, err := o.Pow(exp)
	if err != nil {
		return nil, err
	}

	n2, err := e.Mul(n1)
	if err != nil {
		return nil, err
	}

	return n2, nil
}
