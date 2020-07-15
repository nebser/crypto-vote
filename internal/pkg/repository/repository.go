package repository

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/pkg/errors"
)

func blocksBucket() []byte {
	return []byte("blocks")
}

func tipKey() []byte {
	return []byte("l")
}

type DB struct {
	db *bolt.DB
}

func New(fileName string) (*DB, error) {
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to initialize db with file %s", fileName)
	}
	return &DB{db: db}, nil
}

func (d DB) GetTip() []byte {
	var tip []byte
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket())
		if b == nil {
			return nil
		}
		tip = b.Get(tipKey())
		return nil
	})
	return tip
}

func (d DB) InitBlockchain(genesis blockchain.Block) ([]byte, error) {
	var tip []byte
	err := d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket(blocksBucket())
		if err != nil {
			return errors.Wrap(err, "Failed to create blocks bucket")
		}
		rawBlock, err := json.Marshal(newBlock(genesis))
		if err != nil {
			return errors.Wrapf(err, "Failed to marshal block %#v", genesis)
		}
		if err := b.Put(genesis.Header.TransactionHash, rawBlock); err != nil {
			return errors.Wrap(err, "Failed to put genesis block")
		}
		if err := b.Put(tipKey(), genesis.Header.TransactionHash); err != nil {
			return errors.Wrap(err, "Failed to update tip")
		}
		tip = genesis.Header.TransactionHash
		return nil
	})

	return tip, err
}
