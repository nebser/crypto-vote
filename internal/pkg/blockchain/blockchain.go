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
	db  *bolt.DB
	Tip []byte
}

type GetBlockchainFn func() (*Blockchain, error)

type InitBlockchainFn func(Block) (*Blockchain, error)

type AddBlockFn func(Blockchain, Block) (*Blockchain, error)

func GetBlockchain(db *bolt.DB) GetBlockchainFn {
	return func() (*Blockchain, error) {
		tip := getTip(db)
		if tip == nil {
			return nil, nil
		}
		return &Blockchain{Tip: tip}, nil
	}
}

func InitBlockchain(db *bolt.DB) InitBlockchainFn {
	return func(genesis Block) (*Blockchain, error) {
		tip, err := initBlockchain(db, genesis)
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

func NewBlockchain(db *bolt.DB) Blockchain {
	tip := getTip(db)
	return Blockchain{Tip: tip, db: db}
}

func (b Blockchain) SetGenesis(genesis Block) (Blockchain, error) {
	tip, err := initBlockchain(b.db, genesis)
	if err != nil {
		return b, errors.Wrap(err, "Failed to initialize blockchain")
	}
	return Blockchain{
		Tip: tip,
		db:  b.db,
	}, nil
}

func (b Blockchain) AddBlock(block Block) (Blockchain, error) {
	tip, err := addBlock(b.db, block)
	if err != nil {
		return Blockchain{}, errors.Wrap(err, "Failed to add block")
	}
	return Blockchain{
		Tip: tip,
		db:  b.db,
	}, nil
}
