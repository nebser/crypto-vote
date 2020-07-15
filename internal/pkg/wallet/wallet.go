package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ripemd160"
)

const (
	version       byte = 0
	addressLength int  = 4
)

type Wallet struct {
	Address    string
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
}

func ExtractAddress(publicKey []byte) (string, error) {
	publicSHA256 := sha256.Sum256(publicKey)
	RIPEMD160Hasher := ripemd160.New()
	if _, err := RIPEMD160Hasher.Write(publicSHA256[:]); err != nil {
		return "", errors.Wrap(err, "Failed to write to RIPEMD")
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	versionedPublicKey := append([]byte{version}, publicRIPEMD160...)
	checksum := getChecksum(versionedPublicKey)

	payload := append(versionedPublicKey, checksum...)
	encoded := base58.Encode(payload)
	return encoded, nil
}

func getChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressLength]
}

func New() (*Wallet, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate private key")
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	address, err := ExtractAddress(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create address")
	}

	return &Wallet{
		PublicKey:  pubKey,
		PrivateKey: *private,
		Address:    address,
	}, nil
}
