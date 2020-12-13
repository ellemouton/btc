package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ellemouton/btc/hdkeys"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)
var (
	xpriv string
	xpub string
	path string
	seed string
	mnemonic string
	password string
)

func main() {
	app := &cli.App{
		Name: "keytool",
		Usage: "HD keychain tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "xpriv",
				Value:       "",
				Usage:       "Extended priv key",
				Destination: &xpriv,
			},
			&cli.StringFlag{
				Name:        "xpub",
				Value:       "",
				Usage:       "Extended pub key",
				Destination: &xpub,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       "0",
				Usage:       "derivation path",
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "seed",
				Value:       "",
				Usage:       "seed (hex)",
				Destination: &seed,
			},
			&cli.StringFlag{
				Name:        "mnemonic",
				Value:       "",
				Usage:       "mnemonic words",
				Destination: &mnemonic,
			},
			&cli.StringFlag{
				Name:        "password",
				Value:       "",
				Usage:       "password",
				Destination: &password,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "generate a new seed and extended key pair",
				Action: genNew,
			},
			{
				Name:    "fromSeed",
				Aliases: []string{"s"},
				Usage:   "derive extended key pair from seed",
				Action: genFromSeed,
			},
			{
				Name:    "pub",
				Aliases: []string{"p"},
				Usage:   "get xpub from xpriv",
				Action: getPub,
			},
			{
				Name:    "child",
				Aliases: []string{"c"},
				Usage:   "derive child given xpriv/xpub and path",
				Action: getChild,
			},
			{
				Name: "fromMnemonic",
				Aliases: []string{"m"},
				Usage: "derive extended key from mnemonic",
				Action: genFromMnemonic,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getChild(_ *cli.Context) error {
	if path == "" {
		path = "m"
	}

	var s string
	if xpriv != "" {
		s = xpriv
	} else if xpub != "" {
		s = xpub
	} else {
		log.Fatal("either 'xpriv' or 'xpub' must be specified")
	}

	key, err := hdkeys.Parse(s)
	if err != nil {
		return err
	}

	child, err := key.ChildFromPath(path)
	if err != nil {
		return err
	}

	if child.IsPrivate {
		pub, err := child.ExtendedPubKey()
		if err != nil {
			return err
		}

		fmt.Println("Priv:\t", base58.Encode(child.Serialize()))
		fmt.Println("Pub:\t", base58.Encode(pub.Serialize()))
	} else {
		fmt.Println("Pub:\t", base58.Encode(child.Serialize()))
	}

	return nil
}

func genFromSeed(_ *cli.Context) error {
	if seed == ""{
		log.Fatal("must provide 'seed' flag")
	}

	s, err := hex.DecodeString(seed)
	if err != nil {
		log.Fatal(err)
	}

	priv, err := hdkeys.ExtendedPrivKeyFromSeed(s)
	if err != nil {
		return err
	}

	pub, err := priv.ExtendedPubKey()
	if err != nil {
		return err
	}

	fmt.Println("Priv:\t", base58.Encode(priv.Serialize()))
	fmt.Println("Pub:\t", base58.Encode(pub.Serialize()))

	return nil
}

func genFromMnemonic(_ *cli.Context) error {
	if mnemonic == ""{
		log.Fatal("must provide 'mnemonic' flag")
	}

	s := bip39.NewSeed(mnemonic, password)

	priv, err := hdkeys.ExtendedPrivKeyFromSeed(s)
	if err != nil {
		return err
	}

	pub, err := priv.ExtendedPubKey()
	if err != nil {
		return err
	}

	fmt.Println("Priv:\t", base58.Encode(priv.Serialize()))
	fmt.Println("Pub:\t", base58.Encode(pub.Serialize()))

	return nil
}

func getPub(_ *cli.Context) error {
	if xpriv == "" {
		log.Fatal("must provide 'xpriv' flag")
	}

	k, err := hdkeys.Parse(xpriv)
	if err != nil {
		return err
	}

	pub, err := k.ExtendedPubKey()
	if err != nil {
		return err
	}

	fmt.Println("Priv:\t", base58.Encode(k.Serialize()))
	fmt.Println("Pub:\t", base58.Encode(pub.Serialize()))

	return nil
}

func genNew(_ *cli.Context) error {
	seed, err := hdkeys.NewSeed()
	if err != nil {
		return err
	}

	priv, err := hdkeys.ExtendedPrivKeyFromSeed(seed)
	if err != nil {
		return err
	}

	pub, err := priv.ExtendedPubKey()
	if err != nil {
		return err
	}

	fmt.Println("Priv:\t", base58.Encode(priv.Serialize()))
	fmt.Println("Pub:\t", base58.Encode(pub.Serialize()))

	return nil
}
