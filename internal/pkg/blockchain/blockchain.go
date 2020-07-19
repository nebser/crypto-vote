package blockchain

import (
	"github.com/pkg/errors"
)

const (
	magicNumber = 0x100
	version     = 0
)

type GetTipFn func() []byte

type InitBlockchainFn func(Block) ([]byte, error)

type AddBlockFn func(Block) ([]byte, error)

type Blockchain struct {
	getTip         GetTipFn
	initBlockchain InitBlockchainFn
	addBlock       AddBlockFn
	Tip            []byte
}

func NewBlockchain(getTip GetTipFn, initBlockchain InitBlockchainFn, addBlock AddBlockFn) Blockchain {
	tip := getTip()
	return Blockchain{
		Tip:            tip,
		getTip:         getTip,
		initBlockchain: initBlockchain,
		addBlock:       addBlock,
	}
}

func (b Blockchain) SetGenesis(genesis Block) (Blockchain, error) {
	tip, err := b.initBlockchain(genesis)
	if err != nil {
		return b, errors.Wrap(err, "Failed to initialize blockchain")
	}
	b.Tip = tip
	return b, nil
}

func (b Blockchain) AddBlock(block Block) (Blockchain, error) {
	tip, err := b.addBlock(block)
	if err != nil {
		return Blockchain{}, errors.Wrap(err, "Failed to add block")
	}
	b.Tip = tip
	return b, nil
}
