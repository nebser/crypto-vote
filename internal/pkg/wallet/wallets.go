package wallet

import "encoding/json"

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
