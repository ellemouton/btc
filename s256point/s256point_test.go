package s256point

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ellemouton/btc/signature"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r, err := G.Mul(N)
	require.NoError(t, err)

	i, err := New(nil, nil)
	require.NoError(t, err)

	require.Equal(t, r, i)
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name        string
		z           string
		r           string
		s           string
		px          string
		py          string
		expectValid bool
	}{
		{
			name:        "1",
			z:           "bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423",
			r:           "37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6",
			s:           "8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec",
			px:          "04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574",
			py:          "82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4",
			expectValid: true,
		},
		{
			name:        "2",
			z:           "ec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60",
			r:           "ac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395",
			s:           "68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4",
			px:          "887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
			py:          "61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34",
			expectValid: true,
		},
		{
			name:        "3",
			z:           "7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
			r:           "eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
			s:           "c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
			px:          "887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
			py:          "61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34",
			expectValid: true,
		},
		{
			name:        "4",
			z:           "9c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
			r:           "eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
			s:           "c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
			px:          "887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
			py:          "61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34",
			expectValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			z, _ := new(big.Int).SetString(test.z, 16)
			r, _ := new(big.Int).SetString(test.r, 16)
			s, _ := new(big.Int).SetString(test.s, 16)
			px, _ := new(big.Int).SetString(test.px, 16)
			py, _ := new(big.Int).SetString(test.py, 16)

			p, err := New(px, py)
			require.NoError(t, err)

			sig := &signature.Signature{R: r, S: s}

			valid, err := p.Verify(z.Bytes(), sig)
			require.NoError(t, err)

			require.Equal(t, test.expectValid, valid)
		})
	}
}

func TestTemp(t *testing.T) {
	px, _ := new(big.Int).SetString("887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c", 16)
	py, _ := new(big.Int).SetString("61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34", 16)

	p, err := New(px, py)
	require.NoError(t, err)

	fmt.Println(p.SecString(true))
	fmt.Println(p.SecString(false))
}
