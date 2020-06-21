package privatekey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	privKey, err := New(big.NewInt(12345))
	require.NoError(t, err)

	hash := []byte("my test message")

	sig, err := privKey.Sign(hash)
	require.NoError(t, err)

	valid, err := privKey.pubKey.Verify(hash, sig)
	require.NoError(t, err)
	require.True(t, valid)
}
