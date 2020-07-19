package repository

import (
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
)

type block struct {
	MagicNumber      int                      `json:"magicNumber"`
	Size             int                      `json:"blockSize"`
	Version          int                      `json:"versionNumber"`
	PrevBlock        []byte                   `json:"prevBlock"`
	TransactionHash  []byte                   `json:"transactionHash"`
	Timestamp        int64                    `json:"timestamp"`
	TransactionCount int                      `json:"transactionCount"`
	Transactions     transaction.Transactions `json:"transactions"`
	Hash             []byte                   `json:"hash"`
}

func (b block) toBlock() blockchain.Block {
	return blockchain.Block{
		Metadata: blockchain.Metadata{
			MagicNumber: b.MagicNumber,
			Size:        b.Size,
		},
		Header: blockchain.Header{
			Hash:            b.Hash,
			Prev:            b.PrevBlock,
			Timestamp:       b.Timestamp,
			TransactionHash: b.TransactionHash,
			Version:         b.Version,
		},
		Body: blockchain.Body{
			Transactions:      b.Transactions,
			TransactionsCount: b.TransactionCount,
		},
	}
}

func newBlock(b blockchain.Block) block {
	return block{
		MagicNumber:      b.Metadata.MagicNumber,
		Size:             b.Metadata.Size,
		Version:          b.Header.Version,
		PrevBlock:        b.Header.Prev,
		TransactionHash:  b.Header.TransactionHash,
		Timestamp:        b.Header.Timestamp,
		TransactionCount: b.Body.TransactionsCount,
		Transactions:     b.Body.Transactions,
		Hash:             b.Header.Hash,
	}
}
