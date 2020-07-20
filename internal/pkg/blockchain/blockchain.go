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
	Height         int
}

func NewBlockchain(getTip GetTipFn, initBlockchain InitBlockchainFn, addBlock AddBlockFn, getBlock GetBlockFn) (*Blockchain, error) {
	tip := getTip()
	height, err := getHeight(tip, getBlock)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		Tip:            tip,
		Height:         height,
		getTip:         getTip,
		initBlockchain: initBlockchain,
		addBlock:       addBlock,
		getBlock:       getBlock,
	}, nil
}

func getHeight(tip []byte, getBlock GetBlockFn) (int, error) {
	result := 0
	for current := tip; current != nil; {
		block, err := getBlock(current)
		if err != nil {
			return 0, errors.Wrapf(err, "Failed to get block %x", block)
		}
		result++
		current = block.Header.Prev
	}
	return result, nil
}

func (b Blockchain) SetGenesis(genesis Block) (Blockchain, error) {
	tip, err := b.initBlockchain(genesis)
	if err != nil {
		return b, errors.Wrap(err, "Failed to initialize blockchain")
	}
	b.Tip = tip
	b.Height = 1
	return b, nil
}

func (b Blockchain) AddBlock(block Block) (Blockchain, error) {
	tip, err := b.addBlock(block)
	if err != nil {
		return Blockchain{}, errors.Wrap(err, "Failed to add block")
	}
	b.Tip = tip
	b.Height++
	return b, nil
}

func (b Blockchain) GetBlock(hash []byte) (*Block, error) {
	return b.getBlock(hash)
}

func (b Blockchain) Print() error {
	fmt.Printf("Block height: %d\n", b.Height)
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
