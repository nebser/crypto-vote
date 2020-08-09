package wallet

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/pkg/errors"
)

type Wallets []Wallet

type dumpable struct {
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"privateKey"`
	Address    string `json:"address"`
}

func (ws Wallets) Serialized() (json.RawMessage, error) {
	dumpables := make([]dumpable, 0, len(ws))
	for _, w := range ws {
		dumpables = append(dumpables, dumpable{
			PublicKey:  w.PublicKey,
			PrivateKey: w.PrivateKey.D.Bytes(),
			Address:    w.Address,
		})
	}
	return json.Marshal(dumpables)
}

func (ws Wallets) Addresses() (addresses []string) {
	for _, w := range ws {
		addresses = append(addresses, w.Address)
	}
	return
}

func ImportMultiple(keyfilesList keyfiles.KeyFilesList) (Wallets, error) {
	result := Wallets{}
	for _, k := range keyfilesList {
		w, err := Import(k)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to import keys")
		}
		result = append(result, *w)
	}
	return result, nil
}
