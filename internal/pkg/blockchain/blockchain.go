package blockchain

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	magicNumber = 0x100
	version     = 0
)

type GetTipFn func() []byte

type InitBlockchainFn func(Block) ([]byte, error)

type AddBlockFn func(Block) ([]byte, error)

type GetBlockFn func(hash []byte) (*Block, error)

type Blockchain struct {
	getTip         GetTipFn
	initBlockchain InitBlockchainFn
	addBlock       AddBlockFn
	getBlock       GetBlockFn
	Tip            []byte
}

func NewBlockchain(getTip GetTipFn, initBlockchain InitBlockchainFn, addBlock AddBlockFn, getBlock GetBlockFn) Blockchain {
	tip := getTip()
	return Blockchain{
		Tip:            tip,
		getTip:         getTip,
		initBlockchain: initBlockchain,
		addBlock:       addBlock,
		getBlock:       getBlock,
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

func (b Blockchain) GetBlock(hash []byte) (*Block, error) {
	return b.getBlock(hash)
}

func (b Blockchain) Print() error {
	return printOne(b.Tip, b)
}

func printOne(hash []byte, chain Blockchain) error {
	if hash == nil {
		return nil
	}
	block, err := chain.GetBlock(hash)
	if err != nil {
		return errors.Wrapf(err, "Failed to fetch block %s", hash)
	}
	if block == nil {
		return errors.Errorf("Block with hash %s does not exist", hash)
	}
	fmt.Printf("%s", *block)
	return printOne(block.Header.Prev, chain)
}
