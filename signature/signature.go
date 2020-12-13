package signature

import (
	"encoding/hex"
	"errors"
	"math/big"
)

type Signature struct {
	Rx *big.Int
	S  *big.Int
}

func New(r, s *big.Int) *Signature {
	return &Signature{
		Rx: r,
		S:  s,
	}
}

/*
- 0x30 marker byte
- encode len of sig (usually 0x44 or 0x45)
- 0x02 marker byte
- r: as big endian. but prepend with 0x00 if r's 1st byte >= 0x80. prepend resulting len to r.
- 0x02 marker byte
- s: as big endian. but prepend with 0x00 if s's 1st byte >= 0x80. prepend resulting len to s.
*/
func (s *Signature) Der() []byte {
	rbin := s.Rx.Bytes()
	if (rbin[0] & 0x80) >= 1 {
		rbin = append([]byte{0}, rbin...)
	}
	rbin = append([]byte{0x02, byte(len(rbin))}, rbin...)

	sbin := s.S.Bytes()
	if (sbin[0] & 0x80) >= 1 {
		sbin = append([]byte{0}, sbin...)
	}
	sbin = append([]byte{0x02, byte(len(sbin))}, sbin...)

	res := append(rbin, sbin...)

	return append([]byte{0x30, byte(len(res))}, res...)
}

func (s *Signature) DerString() string {
	return hex.EncodeToString(s.Der())
}

func ParseFromString(s string) (*Signature, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return Parse(b)
}

func Parse(sig []byte) (*Signature, error) {
	b := make([]byte, len(sig))
	copy(b, sig)

	if b[0] != 0x30 {
		return nil, errors.New("Bad signature 1")
	}
	b = b[1:]

	length := b[0]
	b = b[1:]
	if int(length)+2 != len(sig) {
		return nil, errors.New("Bad signature Length")
	}

	if b[0] != 0x02 {
		return nil, errors.New("Bad signature 2")
	}
	b = b[1:]

	rlength := b[0]
	b = b[1:]

	r := (&big.Int{}).SetBytes(b[:rlength])
	b = b[rlength:]

	if b[0] != 0x02 {
		return nil, errors.New("Bad signature 3")
	}
	b = b[1:]

	slength := b[0]
	b = b[1:]

	s := (&big.Int{}).SetBytes(b[:slength])
	b = b[slength:]
	if len(b) != 0 {
		return nil, errors.New("Signature too long")
	}

	return &Signature{
		S:  s,
		Rx: r,
	}, nil
}
