package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
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

type AddBlockFn func(Blockchain, Block) (*Blockchain, error)

func GetBlockchain(db *bolt.DB) GetBlockchainFn {
	return func() (*Blockchain, error) {
		tip := getTip(db)
		if tip != nil {
			return &Blockchain{Tip: tip}, nil
		}
		txs := transaction.Transactions{
			transaction.Transaction{
				Outputs: []transaction.Output{
					transaction.Output{Value: 0},
				},
			},
		}
		genesis, err := NewBlock(nil, txs)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create genesis block")
		}
		tip, err = initBlockchain(db, *genesis)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize blockchain")
		}
		return &Blockchain{Tip: tip}, nil
	}
}

func AddBlock(db *bolt.DB) AddBlockFn {
	return func(blockchain Blockchain, block Block) (*Blockchain, error) {
		tip, err := addBlock(db, block)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to add block to db")
		}
		blockchain.Tip = tip
		return &blockchain, nil
	}
}
