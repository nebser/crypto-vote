package repository

import "github.com/nebser/crypto-vote/internal/pkg/blockchain"

type block struct {
	MagicNumber      int      `json:"magicNumber"`
	Size             int      `json:"blockSize"`
	Version          int      `json:"versionNumber"`
	PrevBlock        []byte   `json:"prevBlock"`
	TransactionHash  []byte   `json:"transactionHash"`
	Timestamp        int64    `json:"timestamp"`
	TransactionCount int      `json:"transactionCount"`
	Transactions     []string `json:"transactions"`
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
	}
}
