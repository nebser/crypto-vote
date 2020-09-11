package repository

import (
	"bytes"
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
	Vout          int    `json:"vout"`
}

type utxos []utxo

func utxoBucket() []byte {
	return []byte("utxos")
}

func newUTXO(u transaction.UTXO) utxo {
	return utxo{
		TransactionID: base64.StdEncoding.EncodeToString(u.TransactionID),
		PublicKeyHash: base64.StdEncoding.EncodeToString(u.PublicKeyHash),
		Value:         u.Value,
		Vout:          u.Vout,
	}
}

func (u utxo) toUTXO() transaction.UTXO {
	id, _ := base64.StdEncoding.DecodeString(u.TransactionID)
	publicKeyHash, _ := base64.StdEncoding.DecodeString(u.PublicKeyHash)
	return transaction.UTXO{
		TransactionID: id,
		PublicKeyHash: publicKeyHash,
		Value:         u.Value,
		Vout:          u.Vout,
	}
}

func newUTXOs(ut transaction.UTXOs) utxos {
	result := utxos{}
	for _, u := range ut {
		result = append(result, newUTXO(u))
	}
	return result
}

func (ut utxos) toUTXOs() transaction.UTXOs {
	result := transaction.UTXOs{}
	for _, u := range ut {
		result = append(result, u.toUTXO())
	}
	return result
}

func saveUTXOs(tx *bolt.Tx, utxos transaction.UTXOs) error {
	b := tx.Bucket(utxoBucket())
	if b == nil {
		created, err := tx.CreateBucket(utxoBucket())
		if err != nil {
			return errors.Wrapf(err, "Failed to create bucket %s", utxoBucket())
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

func getUTXOs(tx *bolt.Tx, publicKeyHash []byte) (transaction.UTXOs, error) {
	b := tx.Bucket(utxoBucket())
	if b == nil {
		return nil, nil
	}
	raw := b.Get(publicKeyHash)
	if raw == nil {
		return nil, nil
	}
	var utxos utxos
	if err := json.Unmarshal(raw, &utxos); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal utxos")
	}
	return utxos.toUTXOs(), nil
}

func deleteUTXO(tx *bolt.Tx, utxo transaction.UTXO) error {
	b := tx.Bucket(utxoBucket())
	if b == nil {
		return nil
	}
	utxos, err := getUTXOs(tx, utxo.PublicKeyHash)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve utxo for deletion")
	}
	updated := utxos.Filter(func(u transaction.UTXO) bool {
		return bytes.Compare(utxo.TransactionID, u.TransactionID) != 0
	})
	raw, err := json.Marshal(newUTXOs(updated))
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal utxo %#v", utxos)
	}
	if err := b.Put(utxo.PublicKeyHash, raw); err != nil {
		return errors.Wrapf(err, "Failed to store utxo %#v", utxos)
	}
	return nil
}
