package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/pkg/errors"
)

type Signable interface {
	Signable() ([]byte, error)
}

func Verify(data Signable, signature, publicKey []byte) bool {
	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKey)
	x.SetBytes(publicKey[:(keyLen / 2)])
	y.SetBytes(publicKey[(keyLen / 2):])

	r := big.Int{}
	s := big.Int{}
	signatureLen := len(signature)
	r.SetBytes(signature[:(signatureLen / 2)])
	s.SetBytes(signature[(signatureLen / 2):])

	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
	signable, err := data.Signable()
	if err != nil {
		return false
	}
	return ecdsa.Verify(&pubKey, hash(signable), &r, &s)
}

func Sign(data Signable, privateKey ecdsa.PrivateKey) ([]byte, error) {
	signable, err := data.Signable()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to convert to signable %#v", data)
	}
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, hash(signable))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sign %#v", data)
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func hash(data []byte) []byte {
	hashed := sha256.Sum256(data)
	return hashed[:]
}
