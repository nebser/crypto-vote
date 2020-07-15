package blockchain

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type block struct {
	MagicNumber      int      `json:"magicNumber"`
	Size             int      `json:"blockSize"`
	Version          int      `json:"versionNumber"`
	PrevBlock        []byte   `json:"prevBlock"`
	TransactionHash  []byte   `json:"transactionHash"`
	Timestamp        int64    `json:"timestamp"`
	TransactionCount int      `json:"transactionCount"`
	Transactions     []string `json:"transactions"`
}

func newBlock(b Block) block {
	return block{
		MagicNumber:      b.Metadata.MagicNumber,
		Size:             b.Metadata.Size,
		Version:          b.Header.Version,
		PrevBlock:        b.Header.Prev,
		TransactionHash:  b.Header.TransactionHash,
		Timestamp:        b.Header.Timestamp,
		TransactionCount: b.Body.TransactionsCount,
		Transactions:     b.Body.Transactions,
	}
}

func blocksBucket() []byte {
	return []byte("blocks")
}

func tipKey() []byte {
	return []byte("l")
}

func getTip(db *bolt.DB) []byte {
	var tip []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket())
		if b == nil {
			return nil
		}
		tip = b.Get(tipKey())
		return nil
	})
	return tip
}

func initBlockchain(db *bolt.DB, genesis Block) ([]byte, error) {
	var tip []byte
	err := db.Update(func(tx *bolt.Tx) error {
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
