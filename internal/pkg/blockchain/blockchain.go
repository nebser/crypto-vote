package blockchain

import (
	"fmt"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/pkg/errors"
)

const (
	magicNumber  = 0x100
	version      = 0
	MaxBlockSize = 256
)

type GetTipFn func() []byte

type InitBlockchainFn func(Block) ([]byte, error)

type AddBlockFn func(Block) ([]byte, error)

type AddBlocksFn func(Blocks) ([]byte, error)

type GetBlockFn func(hash []byte) (*Block, error)

type FindBlockFn func(criteria func(Block) bool) (Block, bool, error)

type ForgeBlockFn func(transaction.Transactions) (*Block, error)

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

// func ForgeBlock(getTransactions transaction.GetTransactionsFn, getTip GetTipFn, addBlock AddBlockFn) ForgeBlockFn {
// 	return func() (Block, error) {
// 		transactions, err := getTransactions()
// 		if err != nil {
// 			return Block{}, errors.Wrapf(err, "Failed to get transactions, error %s", err)
// 		}
// 		block, err := NewBlock(getTip(), transactions)
// 		if err != nil {
// 			return Block{}, errors.Wrapf(err, "Failed to create block out of transactions %s", transactions)
// 		}
// 		if _, err := addBlock(*block); err != nil {
// 			return Block{}, errors.Wrapf(err, "Failed to add block %s to blockchain", block)
// 		}
// 		return *block, nil
// 	}
// }

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
