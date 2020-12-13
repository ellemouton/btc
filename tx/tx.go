package tx

import "encoding/hex"

type Tx struct {
	Version int64
	// Inputs
	// Outputs
	// Locktime
	IsTestnet bool
}

func (tx *Tx) Hash() ([]byte, error) {
	return nil, nil
}

func (tx *Tx) ID() (string, error) {
	h, err := tx.Hash()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h), nil
}

func (tx *Tx) Serialize() ([]byte, error) {
	return nil, nil
}

func Parse(b []byte) (*Tx, error) {
	return nil, nil
}

func ParseString(s string) (*Tx, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return Parse(b)
}
