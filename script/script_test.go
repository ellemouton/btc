package script

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	b, err := hex.DecodeString("6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")
	require.NoError(t, err)

	s := Parse(b)
	require.Len(t, s, 2)
	require.Equal(t, "304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601", hex.EncodeToString(s[0].data))
	require.Equal(t, "035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937", hex.EncodeToString(s[1].data))

	serBytes, err := s.Serialize()
	require.NoError(t, err)
	require.Equal(t, b, serBytes)
}
