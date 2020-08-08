package helpers

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func DoubleSha256(b []byte) []byte {
	h1 := sha256.New()
	h1.Write(b)

	h2 := sha256.New()
	h2.Write(h1.Sum(nil))

	return h2.Sum(nil)
}

func Hash160(b []byte) []byte {
	h256 := sha256.New()
	h256.Write(b)

	rip160 := ripemd160.New()
	rip160.Write(h256.Sum(nil))

	return rip160.Sum(nil)
}