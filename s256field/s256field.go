package s256field

import (
	"math/big"

	"github.com/ellemouton/btc/fieldelement"
)

var P *big.Int

const p = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"

type S256Field struct {
	fieldelement.FieldElement
}

func New(n *big.Int) (fieldelement.FieldElement, error) {
	f, err := fieldelement.New(n, P)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *S256Field) Sqrt() (*S256Field, error) {
	e := &big.Int{}
	e.Add(P, big.NewInt(1))
	e.Div(e, big.NewInt(4))

	fe, err := s.Pow(e)
	if err != nil {
		return nil, err
	}

	return &S256Field{fe}, nil
}

func init() {
	pVal, ok := new(big.Int).SetString(p, 16)
	if !ok {
		panic("invalid hex: " + p)
	}
	P = pVal
}
