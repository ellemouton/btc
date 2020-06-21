package point

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ellemouton/btc/fieldelement"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	a, err := fieldelement.New(big.NewInt(0), big.NewInt(223))
	require.NoError(t, err)

	b, err := fieldelement.New(big.NewInt(7), big.NewInt(223))
	require.NoError(t, err)

	x, err := fieldelement.New(big.NewInt(192), big.NewInt(223))
	require.NoError(t, err)

	y, err := fieldelement.New(big.NewInt(105), big.NewInt(223))
	require.NoError(t, err)

	p1, err := New(x, y, a, b)
	require.NoError(t, err)

	p2, err := New(x, y, a, b)
	require.NoError(t, err)
	require.Equal(t, p1, p2)

	c, err := fieldelement.New(big.NewInt(6), big.NewInt(223))
	require.NoError(t, err)

	_, err = New(x, y, a, c)
	require.Error(t, err)
}

func TestCopy(t *testing.T) {
	a, err := fieldelement.New(big.NewInt(0), big.NewInt(223))
	require.NoError(t, err)

	b, err := fieldelement.New(big.NewInt(7), big.NewInt(223))
	require.NoError(t, err)

	x, err := fieldelement.New(big.NewInt(192), big.NewInt(223))
	require.NoError(t, err)

	y, err := fieldelement.New(big.NewInt(105), big.NewInt(223))
	require.NoError(t, err)

	p1, err := New(x, y, a, b)
	require.NoError(t, err)

	p2 := p1.Copy()

	require.Equal(t, p1, p2)
	require.False(t, p1 == p2)
}

type point struct {
	a int64
	b int64
	x int64
	y int64
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name  string
		prime int64
		p1    point
		p2    point
		p3    point
	}{
		{
			name:  "1",
			prime: 223,
			p1:    point{a: 0, b: 7, x: 192, y: 105},
			p2:    point{a: 0, b: 7, x: 17, y: 56},
			p3:    point{a: 0, b: 7, x: 170, y: 142},
		},
		{
			name:  "2",
			prime: 223,
			p1:    point{a: 0, b: 7, x: 47, y: 71},
			p2:    point{a: 0, b: 7, x: 117, y: 141},
			p3:    point{a: 0, b: 7, x: 60, y: 139},
		},
		{
			name:  "3",
			prime: 223,
			p1:    point{a: 0, b: 7, x: 143, y: 98},
			p2:    point{a: 0, b: 7, x: 76, y: 66},
			p3:    point{a: 0, b: 7, x: 47, y: 71},
		},
		{
			name:  "4",
			prime: 223,
			p1:    point{a: 0, b: 7, x: 192, y: 105},
			p2:    point{a: 0, b: 7, x: 192, y: 105},
			p3:    point{a: 0, b: 7, x: 49, y: 71},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a, err := fieldelement.New(big.NewInt(test.p1.a), big.NewInt(test.prime))
			require.NoError(t, err)

			b, err := fieldelement.New(big.NewInt(test.p1.b), big.NewInt(test.prime))
			require.NoError(t, err)

			x1, err := fieldelement.New(big.NewInt(test.p1.x), big.NewInt(test.prime))
			require.NoError(t, err)

			y1, err := fieldelement.New(big.NewInt(test.p1.y), big.NewInt(test.prime))
			require.NoError(t, err)

			x2, err := fieldelement.New(big.NewInt(test.p2.x), big.NewInt(test.prime))
			require.NoError(t, err)

			y2, err := fieldelement.New(big.NewInt(test.p2.y), big.NewInt(test.prime))
			require.NoError(t, err)

			p1, err := New(x1, y1, a, b)
			require.NoError(t, err)

			p2, err := New(x2, y2, a, b)
			require.NoError(t, err)

			p3, err := p1.Add(p2)
			require.NoError(t, err)
			require.Equal(t, a, p3.A)
			require.Equal(t, b, p3.B)

			x, err := fieldelement.New(big.NewInt(test.p3.x), big.NewInt(test.prime))
			require.NoError(t, err)

			y, err := fieldelement.New(big.NewInt(test.p3.y), big.NewInt(test.prime))
			require.NoError(t, err)

			require.Equal(t, x, p3.X)
			require.Equal(t, y, p3.Y)

		})
	}
}

func TestAdd2(t *testing.T) {
	a, err := fieldelement.New(big.NewInt(0), big.NewInt(223))
	require.NoError(t, err)

	b, err := fieldelement.New(big.NewInt(7), big.NewInt(223))
	require.NoError(t, err)

	x2, err := fieldelement.New(big.NewInt(17), big.NewInt(223))
	require.NoError(t, err)

	y2, err := fieldelement.New(big.NewInt(56), big.NewInt(223))
	require.NoError(t, err)

	p1, err := New(nil, nil, a, b)
	require.NoError(t, err)

	p2, err := New(x2, y2, a, b)
	require.NoError(t, err)

	p3, err := p1.Add(p2)
	require.NoError(t, err)
	require.Equal(t, a, p3.A)
	require.Equal(t, b, p3.B)

	require.Equal(t, p2.X, p3.X)
	require.Equal(t, p2.Y, p3.Y)
}

func TestAdd3(t *testing.T) {
	a, err := fieldelement.New(big.NewInt(0), big.NewInt(223))
	require.NoError(t, err)

	b, err := fieldelement.New(big.NewInt(7), big.NewInt(223))
	require.NoError(t, err)

	x1, err := fieldelement.New(big.NewInt(17), big.NewInt(223))
	require.NoError(t, err)

	y1, err := fieldelement.New(big.NewInt(167), big.NewInt(223))
	require.NoError(t, err)

	x2, err := fieldelement.New(big.NewInt(17), big.NewInt(223))
	require.NoError(t, err)

	y2, err := fieldelement.New(big.NewInt(56), big.NewInt(223))
	require.NoError(t, err)

	p1, err := New(x1, y1, a, b)
	require.NoError(t, err)

	p2, err := New(x2, y2, a, b)
	require.NoError(t, err)

	p3, err := p1.Add(p2)
	require.NoError(t, err)
	require.Equal(t, a, p3.A)
	require.Equal(t, b, p3.B)

	require.True(t, reflect.ValueOf(p3.X).IsNil())
	require.True(t, reflect.ValueOf(p3.Y).IsNil())
}

func TestMul(t *testing.T) {
	tests := []struct {
		name  string
		prime int64
		s     int64
		p1    point
		p2    point
	}{
		{
			name:  "1",
			prime: 223,
			s:     2,
			p1:    point{a: 0, b: 7, x: 192, y: 105},
			p2:    point{a: 0, b: 7, x: 49, y: 71},
		},

		{
			name:  "2",
			prime: 223,
			s:     2,
			p1:    point{a: 0, b: 7, x: 143, y: 98},
			p2:    point{a: 0, b: 7, x: 64, y: 168},
		},
		{
			name:  "3",
			prime: 223,
			s:     2,
			p1:    point{a: 0, b: 7, x: 47, y: 71},
			p2:    point{a: 0, b: 7, x: 36, y: 111},
		},
		{
			name:  "4",
			prime: 223,
			s:     4,
			p1:    point{a: 0, b: 7, x: 47, y: 71},
			p2:    point{a: 0, b: 7, x: 194, y: 51},
		},
		{
			name:  "5",
			prime: 223,
			s:     8,
			p1:    point{a: 0, b: 7, x: 47, y: 71},
			p2:    point{a: 0, b: 7, x: 116, y: 55},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a, err := fieldelement.New(big.NewInt(test.p1.a), big.NewInt(test.prime))
			require.NoError(t, err)

			b, err := fieldelement.New(big.NewInt(test.p1.b), big.NewInt(test.prime))
			require.NoError(t, err)

			x1, err := fieldelement.New(big.NewInt(test.p1.x), big.NewInt(test.prime))
			require.NoError(t, err)

			y1, err := fieldelement.New(big.NewInt(test.p1.y), big.NewInt(test.prime))
			require.NoError(t, err)

			x2, err := fieldelement.New(big.NewInt(test.p2.x), big.NewInt(test.prime))
			require.NoError(t, err)

			y2, err := fieldelement.New(big.NewInt(test.p2.y), big.NewInt(test.prime))
			require.NoError(t, err)

			p1, err := New(x1, y1, a, b)
			require.NoError(t, err)

			p2, err := New(x2, y2, a, b)
			require.NoError(t, err)

			p3, err := p1.Mul(big.NewInt(test.s))
			require.NoError(t, err)

			require.Equal(t, p2, p3)
		})
	}
}

func TestMul2(t *testing.T) {
	a, err := fieldelement.New(big.NewInt(0), big.NewInt(223))
	require.NoError(t, err)

	b, err := fieldelement.New(big.NewInt(7), big.NewInt(223))
	require.NoError(t, err)

	x1, err := fieldelement.New(big.NewInt(47), big.NewInt(223))
	require.NoError(t, err)

	y1, err := fieldelement.New(big.NewInt(71), big.NewInt(223))
	require.NoError(t, err)

	p1, err := New(x1, y1, a, b)
	require.NoError(t, err)

	p3, err := New(nil, nil, a, b)

	p2, err := p1.Mul(big.NewInt(21))
	require.NoError(t, err)
	require.Equal(t, p3, p2)
}
