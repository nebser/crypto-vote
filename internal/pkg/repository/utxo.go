package repository

import (
	"encoding/base64"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/pkg/errors"
)

type utxo struct {
	PublicKeyHash string `json:"publicKeyHash"`
	TransactionID string `json:"transactionId"`
	Value         int    `json:"value"`
}

func UTXOBucket() []byte {
	return []byte("utxos")
}

func newUTXO(u transaction.UTXO) utxo {
	return utxo{
		TransactionID: base64.StdEncoding.EncodeToString(u.TransactionID),
		PublicKeyHash: base64.StdEncoding.EncodeToString(u.PublicKeyHash),
		Value:         u.Value,
	}
}

func (u utxo) toUTXO() transaction.UTXO {
	id, _ := base64.StdEncoding.DecodeString(u.TransactionID)
	publicKeyHash, _ := base64.StdEncoding.DecodeString(u.PublicKeyHash)
	return transaction.UTXO{
		TransactionID: id,
		PublicKeyHash: publicKeyHash,
		Value:         u.Value,
	}
}

func saveUTXOS(tx *bolt.Tx, utxos []transaction.UTXO) error {
	b := tx.Bucket(UTXOBucket())
	if b == nil {
		created, err := tx.CreateBucket(UTXOBucket())
		if err != nil {
			return errors.Wrapf(err, "Failed to create bucket %s", UTXOBucket())
		}
		b = created
	}
	for _, u := range utxos {
		var saved []utxo
		raw := b.Get(u.PublicKeyHash)
		if raw != nil {
			if err := json.Unmarshal(raw, &saved); err != nil {
				return errors.Wrap(err, "Failed to unmarshal into utxo array")
			}
		}
		saved = append(saved, newUTXO(u))
		serialized, err := json.Marshal(saved)
		if err != nil {
			return errors.Wrapf(err, "Failed to serialize %#v", saved)
		}
		if err := b.Put(u.PublicKeyHash, serialized); err != nil {
			return errors.Wrapf(err, "Failed to save utxo set for %x", u.PublicKeyHash)
		}
	}
	return nil
}

func saveUTXOs(transactions transaction.Transactions) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		for _, transaction := range transactions {
			utxos := transaction.UTXOs()
			b := tx.Bucket(UTXOBucket())
			if b == nil {
				created, err := tx.CreateBucket(UTXOBucket())
				if err != nil {
					return errors.Wrapf(err, "Failed to create bucket %s", UTXOBucket())
				}
				b = created
			}
			for _, u := range utxos {
				var saved []utxo
				raw := b.Get(u.PublicKeyHash)
				if raw != nil {
					if err := json.Unmarshal(raw, &saved); err != nil {
						return errors.Wrap(err, "Failed to unmarshal into utxo array")
					}
				}
				saved = append(saved, newUTXO(u))
				serialized, err := json.Marshal(saved)
				if err != nil {
					return errors.Wrapf(err, "Failed to serialize %#v", saved)
				}
				if err := b.Put(u.PublicKeyHash, serialized); err != nil {
					return errors.Wrapf(err, "Failed to save utxo set for %x", u.PublicKeyHash)
				}
			}
		}
		return nil
	}
}
