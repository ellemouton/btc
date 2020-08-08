package hdkeys

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ellemouton/btc/helpers"
	"github.com/ellemouton/btc/privatekey"
	"github.com/ellemouton/btc/s256point"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	seedSize = 64
)

var (
	XprivVersion, _ = hex.DecodeString("0488ade4")
	XpubVersion, _ = hex.DecodeString("0488b21e")
)

func ExtendedPrivKeyFromSeed(s []byte) (*ExtendedKey, error) {
	hmac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	_, err := hmac.Write(s)
	if err != nil {
		return nil, err
	}

	I := hmac.Sum(nil)
	
	key := I[:32]
	chaincode := I[32:]

	if err = validatePrivateKey(key); err != nil {
		return nil, err
	}
	
	return &ExtendedKey{
		Version: XprivVersion,
		Key: key,
		ChainCode:  chaincode,
		Depth: 0,
		FingerPrint: []byte{0x0, 0x0, 0x0, 0x0},
		Index: 0,
		IsPrivate: true,
	}, nil
}

func (priv *ExtendedKey) ExtendedPubKey() (*ExtendedKey, error) {
	if !priv.IsPrivate {
		return nil, errors.New("key is already a public key")
	}

	privKey, err := privatekey.New(new(big.Int).SetBytes(priv.Key))
	if err != nil {
		return nil, err
	}

	pubKey := privKey.PubKey.Sec(true)

	return &ExtendedKey{
		Version: XpubVersion,
		Key: pubKey,
		ChainCode:  priv.ChainCode,
		Depth: priv.Depth,
		FingerPrint: priv.FingerPrint,
		Index: priv.Index,
		IsPrivate: false,
	}, nil
}

func (ext *ExtendedKey) Child(i uint32) (*ExtendedKey, error) {
	if !ext.IsPrivate && i >= uint32(0x80000000) {
		return nil, errors.New("cant derive a hardened child from a public key")
	}

	index := uint32Bytes(i)

	var data []byte
	if i >= uint32(0x80000000) {
		// hardened. so private key
		data = append([]byte{0x0}, ext.Key...)
	} else {
		// normal. so we add the pub key bytes
		if ext.IsPrivate {
			privKey, err := privatekey.New(new(big.Int).SetBytes(ext.Key))
			if err != nil {
				return nil, err
			}
			data = privKey.PubKey.Sec(true)
		} else {
			data = ext.Key
		}
	}
	data = append(data, index...)

	hmac := hmac.New(sha512.New, ext.ChainCode)
	_, err := hmac.Write(data)
	if err != nil {
		return nil, err
	}
	constant := hmac.Sum(nil)

	child := &ExtendedKey{
		Depth: ext.Depth + 1,
		ChainCode: constant[32:],
		IsPrivate: ext.IsPrivate,
		Index: int64(i),
	}

	if ext.IsPrivate {
		child.Version = XprivVersion
		child.Key = addPrivKeys(constant[:32], ext.Key)

		privKey, err := privatekey.New(new(big.Int).SetBytes(ext.Key))
		if err != nil {
			return nil, err
		}
		child.FingerPrint = helpers.Hash160(privKey.PubKey.Sec(true))[:4]

	} else {
		child.Version = XpubVersion
		child.FingerPrint = helpers.Hash160(ext.Key)[:4]

		privKey, err := privatekey.New(new(big.Int).SetBytes(constant[:32]))
		if err != nil {
			return nil, err
		}

		pubKey := privKey.PubKey.Sec(true)

		p1, err := s256point.Parse(pubKey)
		if err != nil {
			return nil, err
		}

		p2, err := s256point.Parse(ext.Key)
		if err != nil {
			return nil, err
		}

		p3, err := p1.Add(p2)
		if err != nil {
			return nil, err
		}

		child.Key = (&s256point.S256Point{p3}).Sec(true)
	}

	return child, nil
}

func (ext *ExtendedKey) ChildFromPath(path string) (*ExtendedKey, error) {
	pathArr, err := getPath(path)
	if err != nil {
		return nil, err
	}

	key, err := ext.Clone()
	if err != nil {
		return nil, err
	}

	for _, p := range pathArr {
		key, err = key.Child(p)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}

func addPrivKeys(key1 []byte, key2 []byte) []byte {
	var key1Int big.Int
	var key2Int big.Int
	key1Int.SetBytes(key1)
	key2Int.SetBytes(key2)

	key1Int.Add(&key1Int, &key2Int)
	key1Int.Mod(&key1Int, s256point.N)

	b := key1Int.Bytes()
	if len(b) < 32 {
		extra := make([]byte, 32-len(b))
		b = append(extra, b...)
	}
	return b
}


func uint32Bytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}

func validatePrivateKey(key []byte) error {
	if fmt.Sprintf("%x", key) == "0000000000000000000000000000000000000000000000000000000000000000" || //if the key is zero
		bytes.Compare(key, s256point.N.Bytes()) >= 0 || //or is outside of the curve
		len(key) != 32 { //or is too short
		return errors.New("invalid priv key bytes")
	}

	return nil
}

func NewSeed() ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, seedSize)
	_, err := rand.Read(token)
	return token, err
}

func getPath(path string) ([]uint32, error) {
	p := strings.Split(path, "/")

	var final []uint32

	if len(p) == 1 && p[0] == "m" {
		return nil, nil
	}

	if p[0] != "m" {
		return nil, errors.New("path must start with 'm'")
	}

	for _, v := range p[1:] {
		var i uint32
		if strings.HasSuffix(v, "'") {
			// hardened child
			num, err := strconv.Atoi(v[:len(v)-1])
			if err != nil {
				return nil, err
			}
			i = uint32(1<<31 + num)
		} else {
			// normal child
			num, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			i = uint32(num)
		}
		final = append(final, i)
	}

	return final, nil
}
