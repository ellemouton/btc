package s256field

import (
	"math/big"

	"github.com/ellemouton/btc/fieldelement"
)

var P *big.Int

const p = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"

type S256Field interface {
	fieldelement.FieldElement
}

func New(n *big.Int) (S256Field, error) {
	f, err := fieldelement.New(n, P)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func Sqrt(s S256Field) (S256Field, error) {
	e := &big.Int{}
	e.Add(P, big.NewInt(1))
	e.Div(e, big.NewInt(4))

	fe, err := s.Pow(e)
	if err != nil {
		return nil, err
	}

	return fe, nil
}

func init() {
	pVal, ok := new(big.Int).SetString(p, 16)
	if !ok {
		panic("invalid hex: " + p)
	}
	P = pVal
}
