package privatekey

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"

	"github.com/ellemouton/btc/s256point"
	"github.com/ellemouton/btc/signature"
)

type PrivateKey struct {
	secret *big.Int
	PubKey *s256point.S256Point
}

func New(s *big.Int) (*PrivateKey, error) {
	p, err := s256point.G.Mul(s)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		secret: s,
		PubKey: p.(*s256point.S256Point),
	}, nil
}

func (p *PrivateKey) Hex() string {
	return fmt.Sprintf("%064x", p.secret)
}

func (p *PrivateKey) Sign(hash []byte) (*signature.Signature, error) {
	k := p.DeterministicK(hash)

	z := hashToInt(hash, s256point.N)

	R, err := s256point.G.Mul(k)
	if err != nil {
		return nil, err
	}

	r := R.GetX().GetNum()

	exp := &big.Int{}
	exp.Sub(s256point.N, big.NewInt(2))

	kInv := &big.Int{}
	kInv.Exp(k, exp, s256point.N)

	s := &big.Int{}
	s.Mul(r, p.secret)
	s.Add(s, z)
	s.Mul(s, kInv)
	s.Mod(s, s256point.N)

	halfN := &big.Int{}
	halfN.Div(s256point.N, big.NewInt(2))

	if s.Cmp(halfN) > 0 {
		s.Sub(s256point.N, s)
	}

	return &signature.Signature{Rx: r, S: s}, nil
}

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

func (p *PrivateKey) DeterministicK(hash []byte) *big.Int {
	n := s256point.N
	e := p.secret
	alg := sha256.New

	qlen := n.BitLen()
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(e, rolen), bits2octets(hash, n, rolen)...)

	// Step B
	v := bytes.Repeat(oneInitializer, holen)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, holen)

	// Step D
	k = mac(alg, k, append(append(v, 0x00), bx...))

	// Step E
	v = mac(alg, k, v)

	// Step F
	k = mac(alg, k, append(append(v, 0x01), bx...))

	// Step G
	v = mac(alg, k, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t)*8 < qlen {
			v = mac(alg, k, v)
			t = append(t, v...)
		}

		// Step H3
		secret := hashToInt(t, n)
		if secret.Cmp(one) >= 0 && secret.Cmp(n) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00))
		v = mac(alg, k, v)
	}
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// left pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, n *big.Int, rolen int) []byte {
	z1 := hashToInt(in, n)
	z2 := new(big.Int).Sub(z1, n)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, n *big.Int) *big.Int {
	orderBits := n.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
