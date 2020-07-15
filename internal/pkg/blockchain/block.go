package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/pkg/errors"
)

type Metadata struct {
	MagicNumber int
	Size        int
}

type Header struct {
	Version         int
	Prev            []byte
	TransactionHash []byte
	Hash            []byte
	Timestamp       int64
}

type Body struct {
	TransactionsCount int
	Transactions      transaction.Transactions
}

type Block struct {
	Metadata Metadata
	Header   Header
	Body     Body
}

func NewBlock(previousBlock []byte, transactions transaction.Transactions) (*Block, error) {
	header := Header{
		Prev:            previousBlock,
		TransactionHash: transactions.Hash(),
		Timestamp:       time.Now().Unix(),
	}
	timestampBytes, err := intToHex(header.Timestamp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to convert timestamp %d to byte array", header.Timestamp)
	}
	hashable := bytes.Join(
		[][]byte{
			header.Prev,
			header.TransactionHash,
			timestampBytes,
		},
		[]byte{},
	)
	hash := sha256.Sum256(hashable)
	header.Hash = hash[:]
	return &Block{
		Header: header,
		Metadata: Metadata{
			MagicNumber: magicNumber,
		},
		Body: Body{
			Transactions:      transactions,
			TransactionsCount: len(transactions),
		},
	}, nil
}

func intToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		return nil, errors.Wrapf(err, "Failed to convert int to byte array %d", num)
	}
	return buff.Bytes(), nil
}
