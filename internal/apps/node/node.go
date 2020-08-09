package node

import (
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/operations"
	"github.com/pkg/errors"
)

func Initialize(
	getHeight operations.GetHeightFn,
	getMissingBlocks operations.GetMissingBlocksFn,
	getBlock operations.GetBlockFn,
	getTip blockchain.GetTipFn,
	getBlockchainBlock blockchain.GetBlockFn,
	addBlocks blockchain.AddBlocksFn,
) error {
	blockchainHeight, err := getHeight()
	if err != nil {
		return errors.Wrap(err, "Couldn't obtain blockchain height")
	}
	localHeight, err := blockchain.GetHeight(getTip, getBlockchainBlock)
	if err != nil {
		return errors.Wrap(err, "Couldn't obtain local blockchain height")
	}
	if localHeight == blockchainHeight {
		return nil
	}
	tip := getTip()
	blockHashes, err := getMissingBlocks(tip)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve missing blocks since tip %x", tip)
	}
	if len(blockHashes) == 0 {
		return nil
	}
	blocks := blockchain.Blocks{}
	for _, hash := range blockHashes {
		block, err := getBlock(hash)
		if err != nil {
			return errors.Wrapf(err, "Failed to obtain block %x", hash)
		}
		blocks = append(blocks, block)
	}
	if _, err := addBlocks(blocks); err != nil {
		return errors.Wrap(err, "Failed to add blocks during initialization")
	}
	return nil
}
