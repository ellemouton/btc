package hdkeys

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSeed(t *testing.T) {
	s, err := NewSeed()
	require.NoError(t, err)
	require.Len(t, s, 64)
}

func TestExtendedKeySerialize(t *testing.T) {
	tests := []struct {
		seed string
		expectPrivSer string
		expectPubSer string
	} {
		{
			seed: "000102030405060708090a0b0c0d0e0f",
			expectPrivSer: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
			expectPubSer: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
		},
		{
			seed: "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
			expectPrivSer: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
			expectPubSer: "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			seed, _ := hex.DecodeString(test.seed)

			k, err := ExtendedPrivKeyFromSeed(seed)
			require.NoError(t, err)

			ser := base58.Encode(k.Serialize())
			require.Equal(t, test.expectPrivSer, ser)

			p, err := k.ExtendedPubKey()
			require.NoError(t, err)

			ser = base58.Encode(p.Serialize())
			require.Equal(t, test.expectPubSer, ser)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []string {
		"xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
		"xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
	}

	for _, test := range tests{
		k, err := Parse(test)
		require.NoError(t, err)

		require.Equal(t, test, base58.Encode(k.Serialize()))
	}
}

func TestChildDerivation(t *testing.T) {
	tests := []struct {
		wif string
		childIndex uint32
		expectChildSer string
	} {
		{
			wif: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
			childIndex: 0,
			expectChildSer: "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
		},
		{
			wif: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
			childIndex: 2147483648,
			expectChildSer: "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7",
		},
		{
			wif: "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
			childIndex: 0,
			expectChildSer: "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
		},
	}

	for _, test := range tests {

		key, err := Parse(test.wif)
		require.NoError(t, err)

		c, err := key.Child(test.childIndex)
		require.NoError(t,err)
		require.Equal(t, test.expectChildSer, base58.Encode(c.Serialize()))
	}
}