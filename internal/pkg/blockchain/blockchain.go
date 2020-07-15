package blockchain

import "github.com/pkg/errors"

const (
	magicNumber = 0x100
	version     = 0
)

type Blockchain struct {
	Tip []byte
}

type GetTipFn func() []byte

type InitBlockchainFn func(Block) ([]byte, error)

func NewBlockchain(getTip GetTipFn, initBlockchain InitBlockchainFn) (*Blockchain, error) {
	tip := getTip()
	if tip != nil {
		return &Blockchain{Tip: tip}, nil
	}
	genesis := Block{
		Metadata: Metadata{
			MagicNumber: magicNumber,
		},
		Header: Header{
			Version:         version,
			TransactionHash: []byte("hash"),
		},
		Body: Body{
			TransactionsCount: 1,
			Transactions:      []string{"transaction"},
		},
	}
	tip, err := initBlockchain(genesis)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize blockchain")
	}
	return &Blockchain{Tip: tip}, nil
}
