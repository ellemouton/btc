package signature

import "math/big"

type Signature struct {
	R *big.Int
	S *big.Int
}
