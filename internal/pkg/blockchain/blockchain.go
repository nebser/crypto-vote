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

type AddBlocksFn func(Blocks) ([]byte, error)

type GetBlockFn func(hash []byte) (*Block, error)

type FindBlockFn func(criteria func(Block) bool) (Block, bool, error)

func GetHeight(getTip GetTipFn, getBlock GetBlockFn) (int, error) {
	result := 0
	for current := getTip(); current != nil; {
		block, err := getBlock(current)
		if err != nil {
			return 0, errors.Wrapf(err, "Failed to get block %x", block)
		}
		result++
		current = block.Header.Prev
	}
	return result, nil
}

func FindBlock(getTip GetTipFn, getBlock GetBlockFn) FindBlockFn {
	return func(criteria func(Block) bool) (Block, bool, error) {
		for current := getTip(); current != nil; {
			block, err := getBlock(current)
			if err != nil {
				return Block{}, false, errors.Wrapf(err, "Failed to get block %x", block)
			}
			if criteria(*block) {
				return *block, true, nil
			}
			current = block.Header.Prev
		}
		return Block{}, false, nil
	}
}

func PrintBlockchain(getTip GetTipFn, getBlock GetBlockFn) error {
	height, err := GetHeight(getTip, getBlock)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch height")
	}
	fmt.Printf("Block height: %d\n", height)
	return printOne(getTip(), getBlock)
}

func printOne(hash []byte, getBlock GetBlockFn) error {
	if hash == nil {
		return nil
	}
	block, err := getBlock(hash)
	if err != nil {
		return errors.Wrapf(err, "Failed to fetch block %s", hash)
	}
	if block == nil {
		return errors.Errorf("Block with hash %s does not exist", hash)
	}
	fmt.Printf("%s", *block)
	return printOne(block.Header.Prev, getBlock)
}
