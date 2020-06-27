package s256point

import (
	"encoding/hex"
	"math/big"
	"reflect"

	"github.com/ellemouton/btc/point"
	"github.com/ellemouton/btc/s256field"
	"github.com/ellemouton/btc/signature"
)

var (
	N *big.Int
	G *S256Point
)

const (
	PubKeyBytesLenCompressed   = 33
	PubKeyBytesLenUncompressed = 65
)

type S256Point struct {
	*point.Point
}

func New(x, y *big.Int) (*S256Point, error) {
	a, err := s256field.New(big.NewInt(0))
	if err != nil {
		return nil, err
	}

	b, err := s256field.New(big.NewInt(7))
	if err != nil {
		return nil, err
	}

	if x == nil {
		p, err := point.New(nil, nil, a, b)
		return &S256Point{p}, err
	}

	xf, err := s256field.New(x)
	if err != nil {
		return nil, err
	}

	yf, err := s256field.New(y)
	if err != nil {
		return nil, err
	}

	p, err := point.New(xf, yf, a, b)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}

func (s *S256Point) Add(o *S256Point) (*S256Point, error) {
	p, err := s.Point.Add(o.Point)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}

func (s *S256Point) Mul(c *big.Int) (*S256Point, error) {
	coef := &big.Int{}
	coef.Mod(c, N)

	p, err := s.Point.Mul(coef)
	if err != nil {
		return nil, err
	}

	return &S256Point{p}, nil
}

const (
	pubkeyCompressedEven byte = 0x2 // y_bit + x coord
	pubkeyCompressedOdd  byte = 0x3 // y_bit + x coord
	pubkeyUncompressed   byte = 0x4 // x coord + y coord
)

// Uncompressed (65 bytes):
// - 0x04
// - x coordinate (32 bytes)
// - y coordinate (32 bytes)
// Compressed (33 bytes):
// - even y: -0x2 + x coord
// - odd y: -0x3 + x coord
func (s *S256Point) Sec(compressed bool) []byte {
	if !compressed {
		b := make([]byte, 0, PubKeyBytesLenUncompressed)
		b = append(b, pubkeyUncompressed)
		b = paddedAppend(32, b, s.X.Num.Bytes())
		return paddedAppend(32, b, s.Y.Num.Bytes())
	}

	b := make([]byte, 0, PubKeyBytesLenCompressed)
	if isOdd(s.Y.Num) {
		b = append(b, pubkeyCompressedOdd)
	} else {
		b = append(b, pubkeyCompressedEven)
	}

	return paddedAppend(32, b, s.X.Num.Bytes())
}

func (s *S256Point) SecString(compressed bool) string {
	return hex.EncodeToString(s.Sec(compressed))
}

func (s *S256Point) Verify(hash []byte, sig *signature.Signature) (bool, error) {
	z := hashToInt(hash, N)

	exp := &big.Int{}
	exp.Sub(N, big.NewInt(2))

	s_inv := &big.Int{}
	s_inv.Exp(sig.S, exp, N)

	u := &big.Int{}
	u.Mul(z, s_inv)
	u.Mod(u, N)

	v := &big.Int{}
	v.Mul(sig.R, s_inv)
	v.Mod(v, N)

	uG, err := G.Mul(u)
	if err != nil {
		return false, err
	}

	vG, err := s.Mul(v)
	if err != nil {
		return false, err
	}

	total, err := uG.Add(vG)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(total.X.Num, sig.R), nil
}

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

func init() {
	nVal, ok := new(big.Int).SetString(n, 16)
	if !ok {
		panic("invalid hex: " + n)
	}
	N = nVal

	gxVal, ok := new(big.Int).SetString(gx, 16)
	if !ok {
		panic("invalid hex: " + gx)
	}

	gyVal, ok := new(big.Int).SetString(gy, 16)
	if !ok {
		panic("invalid hex: " + gy)
	}

	gPoint, err := New(gxVal, gyVal)
	if err != nil {
		panic("error initializing point G")
	}

	G = gPoint
}

const (
	n  = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	gx = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	gy = "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
)

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func isOdd(i *big.Int) bool {
	return i.Bit(0) == 1
}
