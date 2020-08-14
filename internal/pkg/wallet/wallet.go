package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"

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

func (w Wallet) PublicKeyHash() []byte {
	return ExtractPublicKeyHash(w.Address)
}

func (w Wallet) Export(filePrefix string) error {
	encodedPrivateKey, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to encode wallet private key")
	}
	pemEncodedPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: encodedPrivateKey,
	})
	if err := ioutil.WriteFile(filePrefix+".pem", pemEncodedPrivateKey, 0644); err != nil {
		return errors.Wrap(err, "Failed to export private key")
	}

	encodedPublicKey, err := x509.MarshalPKIXPublicKey(&w.PrivateKey.PublicKey)
	if err != nil {
		return errors.Wrapf(err, "Failed to encode public key")
	}
	pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: encodedPublicKey,
	})
	if err := ioutil.WriteFile(filePrefix+"_pub.pem", pemEncodedPublicKey, 0644); err != nil {
		return errors.Wrap(err, "Failed to export public key")
	}

	if err := ioutil.WriteFile(filePrefix+"_address.txt", []byte(w.Address), 0644); err != nil {
		return errors.Wrap(err, "Failed to export address")
	}

	return nil
}

func Import(keyfiles keyfiles.KeyFiles) (*Wallet, error) {
	publicKeyContent, err := ioutil.ReadFile(keyfiles.PublicKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read public key")
	}
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyContent))
	rawPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse public key")
	}
	publicKey, ok := rawPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.Errorf("Failed to case %#v to public key", rawPublicKey)
	}
	publicKey.Curve = elliptic.P256()
	pk := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	address, err := ExtractAddress(pk)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to extract address from %s", pk)
	}

	privateKeyContent, err := ioutil.ReadFile(keyfiles.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read private key")
	}
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyContent))
	privateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse private key")
	}
	privateKey.PublicKey = *publicKey
	return &Wallet{
		PublicKey:  pk,
		PrivateKey: *privateKey,
		Address:    address,
	}, nil
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

func HashedPublicKey(publicKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(publicKey)
	RIPEMD160Hasher := ripemd160.New()
	if _, err := RIPEMD160Hasher.Write(publicSHA256[:]); err != nil {
		return nil, errors.Wrap(err, "Failed to write to RIPEMD")
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160, nil
}

func ExtractPublicKeyHash(address string) []byte {
	decoded := base58.Decode(address)
	return decoded[1 : len(decoded)-addressLength]
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
