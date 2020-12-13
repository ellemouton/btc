package varint

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemp(t *testing.T) {
	tests := []struct {
		name    string
		integer uint64
		hex     string
	}{
		{
			name:    "1",
			integer: 100,
			hex:     "64",
		},
		{
			name:    "2",
			integer: 255,
			hex:     "fdff00",
		},
		{
			name:    "3",
			integer: 555,
			hex:     "fd2b02",
		},
		{
			name:    "4",
			integer: 70015,
			hex:     "fe7f110100",
		},
		{
			name:    "5",
			integer: 18005558675309,
			hex:     "ff6dc7ed3e60100000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := Encode(test.integer)
			require.NoError(t, err)

			h, err := hex.DecodeString(test.hex)
			require.NoError(t, err)

			require.Equal(t, h, b)
			require.Equal(t, test.integer, Read(h))

		})
	}
}
