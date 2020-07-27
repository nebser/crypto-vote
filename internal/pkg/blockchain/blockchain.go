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
}

func NewBlockchain(getTip GetTipFn, initBlockchain InitBlockchainFn, addBlock AddBlockFn, getBlock GetBlockFn) *Blockchain {
	return &Blockchain{
		getTip:         getTip,
		initBlockchain: initBlockchain,
		addBlock:       addBlock,
		getBlock:       getBlock,
	}
}

func (b Blockchain) GetHeight() (int, error) {
	result := 0
	for current := b.getTip(); current != nil; {
		block, err := b.getBlock(current)
		if err != nil {
			return 0, errors.Wrapf(err, "Failed to get block %x", block)
		}
		result++
		current = block.Header.Prev
	}
	return result, nil
}

func (b Blockchain) GetTip() []byte {
	return b.getTip()
}

func (b Blockchain) SetGenesis(genesis Block) error {
	if _, err := b.initBlockchain(genesis); err != nil {
		return errors.Wrap(err, "Failed to initialize blockchain")
	}
	return nil
}

func (b Blockchain) AddBlock(block Block) error {
	if _, err := b.addBlock(block); err != nil {
		return errors.Wrap(err, "Failed to add block")
	}
	return nil
}

func (b Blockchain) GetBlock(hash []byte) (*Block, error) {
	return b.getBlock(hash)
}

func (b Blockchain) Print() error {
	height, err := b.GetHeight()
	if err != nil {
		return errors.Wrap(err, "Failed to fetch height")
	}
	fmt.Printf("Block height: %d\n", height)
	return printOne(b.GetTip(), b)
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
