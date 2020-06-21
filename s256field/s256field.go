package s256field

import (
	"errors"
	"math/big"

	"github.com/ellemouton/btc/fieldelement"
)

const P = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"

func New(n *big.Int) (*fieldelement.FieldElement, error) {
	p, ok := new(big.Int).SetString(P, 16)
	if !ok {
		return nil, errors.New("couldnt convert hex string to big Int")
	}

	f, err := fieldelement.New(n, p)
	if err != nil {
		return nil, err
	}

	return f, nil
}
