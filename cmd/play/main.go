package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ellemouton/btc/fieldelement"
	"github.com/ellemouton/btc/helpers"
	"github.com/ellemouton/btc/s256point"
	"github.com/ellemouton/btc/signature"
)

type info struct {
	address string
	message string
	pubkey  string
	sig     string
}

var inf = []info{
	{
		address: "13yaSqGNDzt1mNW4vrKM9CvD46cTavNabF",
		message: "There is nothing too shocking about this signature",
		pubkey:  "02000000000005689111130e588a12ecda87b2dc5585c6c6ba66a412fa0cce65bc",
		sig:     "ffffffff077b7209dc866fbfa0d2bf67a0c696afffe57a822e2ba90059a2cc7abb998becb4e427650e282522bf9576524665301b807c0e3194cf1e6c795a0cf7",
	},
	{
		address: "1MkanKef93F1iNLKvyijrbbW2k5VaXzDvA",
		message: "Nor this, given a bit of algebra.",
		pubkey:  "03742088316dacf400cea17fdea1dba3bc1e1f58ac0f852fd85545b0ba7ebaee79",
		sig:     "10000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000001000000000000000000000000000",
	},
	{
		address: "13see6qjfupx1YWgRefwEkccZeM8QGTAiJ",
		message: "But can you explain this one?",
		pubkey:  "0200000000000000000000003b78ce563f89a0ed9414f5aa28ad0d96d6795f9c63",
		sig:     "deadbeef2f4a23b0f1954100b76bcb720f7b2ddc4a446dc06b8ffc4e143286e1e441f5f1583f300022ad3d134413a212581bcd36c20c7840d15b4d6b8e8f177f",
	},
}

/*
These are all valid signatures in old Bitcoin Armory style using the message hash function sha256(sha256('Bitcoin Signed Message:\n' + message)).
*/

func pkToAddr(pk string) string {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		log.Fatal(err)
	}

	return base58.CheckEncode(helpers.Hash160(pkBytes), 0x0)
}

var msgTmpl = "Bitcoin Signed Message:\n%s"

func msgHash(msg string) []byte {
	return helpers.DoubleSha256([]byte(fmt.Sprintf(msgTmpl, msg)))
}

type parsed struct {
	z   []byte
	pk  *s256point.S256Point
	sig *signature.Signature
}

func parse(i info) (*parsed, error) {
	var res parsed

	p, err := s256point.ParseFromString(i.pubkey)
	if err != nil {
		return nil, err
	}
	res.pk = &s256point.S256Point{p}

	sig, err := hex.DecodeString(i.sig)
	if err != nil {
		return nil, err
	}

	res.sig = &signature.Signature{
		Rx: new(big.Int).SetBytes(sig[:32]),
		S:  new(big.Int).SetBytes(sig[32:]),
	}

	res.z = msgHash(i.message)

	return &res, nil
}

func parseAll(il []info) ([]parsed, error) {
	var res []parsed

	for _, i := range il {
		p, err := parse(i)
		if err != nil {
			return nil, err
		}

		res = append(res, *p)
	}

	return res, nil
}

func printStuff(p parsed) {
	ry, err := getRy(p.pk, p.z, p.sig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("z:\t", p.z)
	fmt.Println("P_x:\t", p.pk.GetX().GetNum().Bytes())
	fmt.Println("P_y:\t", p.pk.GetY().GetNum().Bytes())
	fmt.Println("rx:\t", p.sig.Rx.Bytes())
	fmt.Println("ry:\t", ry.GetNum().Bytes())
	fmt.Println("s:\t", p.sig.S.Bytes())
	fmt.Println()
	fmt.Println("z:\t", new(big.Int).SetBytes(p.z))
	fmt.Println("P_x:\t", p.pk.GetX().GetNum())
	fmt.Println("P_y:\t", p.pk.GetY().GetNum())
	fmt.Println("rx:\t", p.sig.Rx)
	fmt.Println("ry:\t", ry.GetNum())
	fmt.Println("s:\t", p.sig.S)
	fmt.Println()
	fmt.Println("rx:\t", hex.EncodeToString(p.sig.Rx.Bytes()))
	fmt.Println("ry:\t", hex.EncodeToString(ry.GetNum().Bytes()))
	fmt.Println("s:\t", hex.EncodeToString(p.sig.S.Bytes()))
	fmt.Println()
}

func main() {
	pl, err := parseAll(inf)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range pl {
		printStuff(i)
	}

	var i, e = big.NewInt(2), big.NewInt(52)
	i.Exp(i, e, nil)
	fmt.Println(i.Bytes())

	i, e = big.NewInt(2), big.NewInt(252)
	i.Exp(i, e, nil)
	fmt.Println(i.Bytes())

	i, e = big.NewInt(2), big.NewInt(108)
	i.Exp(i, e, nil)
	fmt.Println(i.Bytes())

	fmt.Println(s256point.G.GetX().GetNum().Bytes())
	fmt.Println(s256point.G.GetY().GetNum().Bytes())
}

func getRy(p *s256point.S256Point, hash []byte, sig *signature.Signature) (fieldelement.FieldElement, error) {
	z := hashToInt(hash, s256point.N)

	exp := &big.Int{}
	exp.Sub(s256point.N, big.NewInt(2))

	s_inv := &big.Int{}
	s_inv.Exp(sig.S, exp, s256point.N)

	u := &big.Int{}
	u.Mul(z, s_inv)
	u.Mod(u, s256point.N)

	v := &big.Int{}
	v.Mul(sig.Rx, s_inv)
	v.Mod(v, s256point.N)

	uG, err := s256point.G.Mul(u)
	if err != nil {
		return nil, err
	}

	vG, err := p.Mul(v)
	if err != nil {
		return nil, err
	}

	total, err := uG.Add(vG)
	if err != nil {
		return nil, err
	}

	return total.GetY(), nil
}

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
