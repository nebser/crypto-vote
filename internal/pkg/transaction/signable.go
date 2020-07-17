package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/pkg/errors"
)

type signable struct {
	Sender    []byte `json:"sender"`
	Recipient []byte `json:"recipient"`
	Value     int    `json:"value"`
}

func sign(data signable, key ecdsa.PrivateKey) ([]byte, error) {
	hash, err := hash(data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to hash signable")
	}
	r, s, err := ecdsa.Sign(rand.Reader, &key, hash)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sign %s", hash)
	}
	return append(r.Bytes(), s.Bytes()...), nil
}
