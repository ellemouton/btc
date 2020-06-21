package s256point

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r, err := G.Mul(N)
	require.NoError(t, err)

	i, err := New(nil, nil)
	require.NoError(t, err)

	require.Equal(t, r, i)
}
