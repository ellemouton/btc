package signature

import "math/big"

type Signature struct {
	R *big.Int
	S *big.Int
}

/*
func (s *Signature) Der() ([]byte, err) {


}
*/
