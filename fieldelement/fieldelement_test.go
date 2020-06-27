package fieldelement

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	f1, err := New(big.NewInt(3), big.NewInt(19))
	require.NoError(t, err)

	f2, err := New(big.NewInt(3), big.NewInt(19))
	require.NoError(t, err)

	f3, err := New(big.NewInt(5), big.NewInt(19))
	require.NoError(t, err)

	f4, err := New(big.NewInt(3), big.NewInt(11))
	require.NoError(t, err)

	_, err = New(big.NewInt(15), big.NewInt(11))
	require.Error(t, err)

	_, err = New(big.NewInt(-3), big.NewInt(11))
	require.Error(t, err)

	require.Equal(t, f1, f2)
	require.NotEqual(t, f1, f3)
	require.NotEqual(t, f1, f4)
}

func TestCopy(t *testing.T) {
	f1, err := New(big.NewInt(3), big.NewInt(19))
	require.NoError(t, err)

	f2 := f1.Copy()

	require.Equal(t, f1, f2)
	require.False(t, f1 == f2)
}

func TestAdd(t *testing.T) {
	f1, err := New(big.NewInt(3), big.NewInt(19))
	require.NoError(t, err)

	f2, err := New(big.NewInt(5), big.NewInt(19))
	require.NoError(t, err)

	f3, err := f1.Add(f2)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(19), f3.GetPrime())
	require.Equal(t, big.NewInt(8), f3.GetNum())

	f4, err := New(big.NewInt(16), big.NewInt(19))
	require.NoError(t, err)

	f5, err := f4.Add(f2)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(19), f5.GetPrime())
	require.Equal(t, big.NewInt(2), f5.GetNum())

	f6, err := New(big.NewInt(3), big.NewInt(11))
	require.NoError(t, err)

	_, err = f6.Add(f1)
	require.Error(t, err)
}

func TestSub(t *testing.T) {
	f1, err := New(big.NewInt(3), big.NewInt(19))
	require.NoError(t, err)

	f2, err := New(big.NewInt(5), big.NewInt(19))
	require.NoError(t, err)

	f3, err := f1.Sub(f2)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(19), f3.GetPrime())
	require.Equal(t, big.NewInt(17), f3.GetNum())

	f4, err := f2.Sub(f1)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(19), f4.GetPrime())
	require.Equal(t, big.NewInt(2), f4.GetNum())

	f5, err := New(big.NewInt(3), big.NewInt(11))
	require.NoError(t, err)

	_, err = f5.Sub(f1)
	require.Error(t, err)
}

func TestMul(t *testing.T) {
	f1, err := New(big.NewInt(24), big.NewInt(31))
	require.NoError(t, err)

	f2, err := New(big.NewInt(19), big.NewInt(31))
	require.NoError(t, err)

	f3, err := New(big.NewInt(22), big.NewInt(31))
	require.NoError(t, err)

	f4, err := f1.Mul(f2)
	require.NoError(t, err)
	require.Equal(t, f3, f4)
}

func TestDiv(t *testing.T) {
	tests := []struct {
		p  int64
		n1 int64
		n2 int64
		n3 int64
	}{
		{p: 19, n1: 2, n2: 7, n3: 3},
		{p: 31, n1: 3, n2: 24, n3: 4},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			f1, err := New(big.NewInt(test.n1), big.NewInt(test.p))
			require.NoError(t, err)

			f2, err := New(big.NewInt(test.n2), big.NewInt(test.p))
			require.NoError(t, err)

			f3, err := New(big.NewInt(test.n3), big.NewInt(test.p))
			require.NoError(t, err)

			f4, err := f1.Div(f2)
			require.NoError(t, err)

			require.Equal(t, f3, f4)
		})
	}
}

func TestDiv2(t *testing.T) {
	f1, err := New(big.NewInt(17), big.NewInt(31))
	require.NoError(t, err)

	f2, err := New(big.NewInt(29), big.NewInt(31))
	require.NoError(t, err)

	f3, err := f1.Pow(big.NewInt(-3))
	require.NoError(t, err)
	require.Equal(t, f2, f3)

	f4, err := New(big.NewInt(4), big.NewInt(31))
	require.NoError(t, err)

	f5, err := New(big.NewInt(11), big.NewInt(31))
	require.NoError(t, err)

	f6, err := New(big.NewInt(13), big.NewInt(31))
	require.NoError(t, err)

	f7, err := f4.Pow(big.NewInt(-4))
	require.NoError(t, err)

	f8, err := f7.Mul(f5)
	require.NoError(t, err)

	require.Equal(t, f6, f8)
}

func TestPow(t *testing.T) {
	tests := []struct {
		p  int64
		c  int64
		n1 int64
		n2 int64
	}{
		{p: 13, c: 3, n1: 3, n2: 1},
		{p: 13, c: -3, n1: 7, n2: 8},
		{p: 31, c: 3, n1: 17, n2: 15},
		{p: 31, c: 3, n1: 17, n2: 15},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			f1, err := New(big.NewInt(test.n1), big.NewInt(test.p))
			require.NoError(t, err)

			f2, err := New(big.NewInt(test.n2), big.NewInt(test.p))
			require.NoError(t, err)

			f3, err := f1.Pow(big.NewInt(test.c))
			require.NoError(t, err)

			require.Equal(t, f2, f3)
		})
	}
}
