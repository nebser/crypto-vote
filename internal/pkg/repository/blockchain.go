package repository

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/pkg/errors"
)

func blocksBucket() []byte {
	return []byte("blocks")
}

func tipKey() []byte {
	return []byte("l")
}

func GetTip(db *bolt.DB) blockchain.GetTipFn {
	return func() []byte {
		var tip []byte
		db.View(func(tx *bolt.Tx) error {
			tip = getTip(tx)
			return nil
		})
		return tip
	}
}

func getTip(tx *bolt.Tx) []byte {
	b := tx.Bucket(blocksBucket())
	if b == nil {
		return nil
	}
	return b.Get(tipKey())
}

func InitBlockchain(db *bolt.DB) blockchain.InitBlockchainFn {
	return func(genesis blockchain.Block) ([]byte, error) {
		var tip []byte
		err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket(blocksBucket())
			if err != nil {
				return errors.Wrap(err, "Failed to create blocks bucket")
			}
			rawBlock, err := json.Marshal(newBlock(genesis))
			if err != nil {
				return errors.Wrapf(err, "Failed to marshal block %#v", genesis)
			}
			if err := b.Put(genesis.Header.Hash, rawBlock); err != nil {
				return errors.Wrap(err, "Failed to put genesis block")
			}
			if err := b.Put(tipKey(), genesis.Header.Hash); err != nil {
				return errors.Wrap(err, "Failed to update tip")
			}
			tip = genesis.Header.Hash
			return nil
		})

		return tip, err
	}
}

func addBlock(tx *bolt.Tx, block blockchain.Block) ([]byte, error) {
	b := tx.Bucket(blocksBucket())
	if b == nil {
		created, err := tx.CreateBucket(blocksBucket())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create blocks bucket")
		}
		b = created
	}
	rawBlock, err := json.Marshal(newBlock(block))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal block %#v", block)
	}
	if err := b.Put(block.Header.Hash, rawBlock); err != nil {
		return nil, errors.Wrapf(err, "Failed to put block %#v", block)
	}
	if err := b.Put(tipKey(), block.Header.Hash); err != nil {
		return nil, errors.Wrap(err, "Failed to update tip")
	}
	return block.Header.Hash, nil
}

func AddBlock(db *bolt.DB) blockchain.AddBlockFn {
	return func(block blockchain.Block) ([]byte, error) {
		var tip []byte
		err := db.Update(func(tx *bolt.Tx) error {
			created, err := addBlockWithUTXO(tx, block)
			if err != nil {
				return errors.Wrapf(err, "Failed to add block %s", block)
			}
			tip = created
			return nil
		})
		return tip, err
	}
}

func addBlockWithUTXO(tx *bolt.Tx, block blockchain.Block) ([]byte, error) {
	tip, err := addBlock(tx, block)
	if err != nil {
		return nil, err
	}
	for _, transaction := range block.Body.Transactions {
		if err := deleteTransaction(tx, transaction); err != nil {
			return nil, err
		}
		if err := saveUTXOs(tx, transaction.UTXOs()); err != nil {
			return nil, err
		}
	}
	return tip, nil
}

func GetBlock(db *bolt.DB) blockchain.GetBlockFn {
	return func(hash []byte) (*blockchain.Block, error) {
		var result *blockchain.Block
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(blocksBucket())
			if b == nil {
				return errors.New("Blocks bucket does not exist")
			}
			rawBlock := b.Get(hash)
			if rawBlock == nil {
				return nil
			}
			var serialized block
			if err := json.Unmarshal(rawBlock, &serialized); err != nil {
				return errors.Wrapf(err, "Failed to unmarshal serialized block %s", rawBlock)
			}
			bl := serialized.toBlock()
			result = &bl
			return nil
		})
		return result, err
	}
}

func ForgeBlock(db *bolt.DB) blockchain.ForgeBlockFn {
	return func(txs transaction.Transactions) (*blockchain.Block, error) {
		var block *blockchain.Block
		err := db.Update(func(tx *bolt.Tx) error {
			var candidates transaction.Transactions
			var invalidTransactions transaction.Transactions
			for _, t := range txs {
				sum, err := getInputSum(tx, t)
				switch {
				case errors.Is(err, transaction.ErrUTXONotFound):
					invalidTransactions = append(invalidTransactions, t)
				case err != nil:
					return errors.Wrapf(err, "Failed to get sum of inputs for transaction %s", t)
				case t.Outputs.Sum() != sum:
					invalidTransactions = append(invalidTransactions, t)
				default:
					candidates = append(candidates, t)
				}
				if len(candidates) == blockchain.MaxBlockSize {
					break
				}
			}
			if err := deleteTransactionsUTXOs(tx, candidates); err != nil {
				return errors.Wrapf(err, "Failed to delete candidate transactions from utxo set %s", candidates)
			}
			if err := deleteTransactions(tx, append(candidates, invalidTransactions...)); err != nil {
				return errors.Wrap(err, "Failed to delete transactions")
			}
			if len(candidates) == 0 {
				return nil
			}
			newBlock, err := blockchain.NewBlock(getTip(tx), candidates)
			if err != nil {
				return errors.Wrap(err, "Failed to set up new block")
			}
			if _, err := addBlockWithUTXO(tx, *newBlock); err != nil {
				return errors.Wrap(err, "Failed to add block to database")
			}
			block = newBlock
			return nil
		})
		return block, err
	}
}
