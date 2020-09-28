package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
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

type Blocks []Block

type VerifyBlockFn func(Block) bool

func (b Block) String() string {
	builder := strings.Builder{}
	builder.WriteString("-----BEGIN BLOCK-----\n")
	builder.WriteString(fmt.Sprintf("Size: %d\n", b.Metadata.Size))
	builder.WriteString(fmt.Sprintf("Hash: %x\n", b.Header.Hash))
	t := time.Unix(b.Header.Timestamp, 0)
	builder.WriteString("Timestamp: ")
	builder.WriteString(t.Format(time.RFC3339))
	builder.WriteString(fmt.Sprintf("\nPrev: %x\n", b.Header.Prev))
	builder.WriteString(b.Body.Transactions.String())
	builder.WriteString("-----END BLOCK-----\n")
	return builder.String()
}

func NewBlock(previousBlock []byte, transactions transaction.Transactions) (*Block, error) {
	transactionsHash := transactions.Hash()
	timestamp := time.Now().Unix()
	blockHash, err := createHash(previousBlock, transactionsHash, timestamp)
	if err != nil {
		return nil, errors.New("Failed to create block hash")
	}
	header := Header{
		Prev:            previousBlock,
		TransactionHash: transactionsHash,
		Timestamp:       timestamp,
		Hash:            blockHash,
	}
	return &Block{
		Header: header,
		Metadata: Metadata{
			MagicNumber: magicNumber,
			Size:        len(transactions),
		},
		Body: Body{
			Transactions:      transactions,
			TransactionsCount: len(transactions),
		},
	}, nil
}

func createHash(previousBlock, transactionsHash []byte, timestamp int64) ([]byte, error) {
	timestampBytes, err := intToHex(timestamp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to convert timestamp %d to byte array", timestamp)
	}
	hashable := bytes.Join(
		[][]byte{
			previousBlock,
			transactionsHash,
			timestampBytes,
		},
		[]byte{},
	)
	hash := sha256.Sum256(hashable)
	return hash[:], nil
}

func intToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		return nil, errors.Wrapf(err, "Failed to convert int to byte array %d", num)
	}
	return buff.Bytes(), nil
}

func VerfiyBlock(verifyTransaction transaction.VerifyTransctionFn, isStakeTransaction transaction.IsStakeTransactionFn) VerifyBlockFn {
	return func(block Block) bool {
		for _, transaction := range block.Body.Transactions {
			if !verifyTransaction(transaction) {
				return false
			}
		}
		if len(block.Body.Transactions) == 0 || !isStakeTransaction(block.Body.Transactions[0]) {
			return false
		}
		transactionHash := block.Body.Transactions.Hash()
		blockHash, err := createHash(block.Header.Prev, transactionHash, block.Header.Timestamp)
		if err != nil {
			return false
		}
		return bytes.Compare(block.Header.Hash, blockHash) == 0
	}
}
