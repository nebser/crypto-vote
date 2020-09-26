package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

type Signer interface {
	Sign(Signable) (string, error)
	Verifier() string
	SignRaw(Signable) ([]byte, error)
}

type walletSigner struct {
	wallet Wallet
}

func (w walletSigner) Sign(signable Signable) (string, error) {
	signature, err := Sign(signable, w.wallet.PrivateKey)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create signature for %#v", signable)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (w walletSigner) SignRaw(signable Signable) ([]byte, error) {
	return Sign(signable, w.wallet.PrivateKey)
}

func (w walletSigner) Verifier() string {
	return base64.StdEncoding.EncodeToString(w.wallet.PublicKey)
}

func NewSigner(wallet Wallet) Signer {
	return walletSigner{wallet: wallet}
}

type VerifierFn func(data Signable, signature, publicKey string) (bool, error)

func VerifySignature(data Signable, signature, publicKey string) (bool, error) {
	rawSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to decode signature %s", signature)
	}
	rawPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to decode public key %s", publicKey)
	}
	return Verify(data, rawSignature, rawPublicKey), nil
}
