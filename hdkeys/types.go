package hdkeys

import (
	"encoding/binary"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ellemouton/btc/helpers"
)

type ExtendedKey struct {
	Version []byte
	Depth int64
	Index int64
	FingerPrint []byte
	Key []byte
	ChainCode []byte
	IsPrivate bool
}

func (ext *ExtendedKey) Clone() (*ExtendedKey, error) {
	temp := ext.Serialize()
	return Parse(base58.Encode(temp))
}

func (ext *ExtendedKey) Serialize() []byte {
	result := make([]byte, 78)

	copy(result[:4], ext.Version)
	copy(result[4:5], []byte{byte(ext.Depth)})
	copy(result[5:9], ext.FingerPrint)
	binary.BigEndian.PutUint32(result[9:13], uint32(ext.Index))
	copy(result[13:45], ext.ChainCode)
	if ext.IsPrivate{
		copy(result[45:78], append([]byte{0x0}, ext.Key...))
	} else {
		copy(result[45:78], ext.Key)
	}

	return append(result, helpers.DoubleSha256(result)[:4]...)
}

func Parse(s string) (*ExtendedKey, error) {
	b := base58.Decode(s)
	if len(b) != 82 {
		return nil, errors.New("incorrect length")
	}

	depth, _ := binary.Varint(b[4:5])
	index, _ := binary.Varint(b[9:13])

	key := b[45:78]
	isPriv := key[0] == byte(0x0)

	if isPriv {
		key = key[1:]
	}

	return &ExtendedKey{
		Version:     b[:4],
		Depth:       depth,
		FingerPrint: b[5:9],
		Index:       index,
		ChainCode:   b[13:45],
		Key:         key,
		IsPrivate:   isPriv,
	}, nil
}

