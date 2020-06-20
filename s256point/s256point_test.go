package s256point

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	g, err := GetG()
	require.NoError(t, err)

	n, err := GetN()
	require.NoError(t, err)

	r, err := g.Mul(n)
	require.NoError(t, err)

	i, err := New(nil, nil)
	require.NoError(t, err)

	require.Equal(t, r, i)
}
