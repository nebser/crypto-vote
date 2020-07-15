package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const (
	magicNumber = 0x100
	version     = 0
)

type Blockchain struct {
	Tip []byte
}

type GetBlockchainFn func() (*Blockchain, error)

func NewBlockchain(db *bolt.DB) GetBlockchainFn {
	return func() (*Blockchain, error) {
		tip := getTip(db)
		if tip != nil {
			return &Blockchain{Tip: tip}, nil
		}
		genesis := Block{
			Metadata: Metadata{
				MagicNumber: magicNumber,
			},
			Header: Header{
				Version:         version,
				TransactionHash: []byte("hash"),
			},
			Body: Body{
				TransactionsCount: 1,
				Transactions:      []string{"transaction"},
			},
		}
		tip, err := initBlockchain(db, genesis)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize blockchain")
		}
		return &Blockchain{Tip: tip}, nil
	}
}
