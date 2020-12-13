package signature

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDER(t *testing.T) {
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)

	p := New(r, s)

	require.Equal(t, "3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", p.DerString())
}

func TestParse(t *testing.T) {
	der := "3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"

	sig, err := ParseFromString(der)
	require.NoError(t, err)

	expectR, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	expectS, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)

	require.True(t, sig.Rx.Cmp(expectR) == 0)
	require.True(t, sig.S.Cmp(expectS) == 0)
}
