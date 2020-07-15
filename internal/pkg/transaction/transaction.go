package transaction

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/pkg/errors"
)

type Transaction struct {
	ID      []byte   `json:"id"`
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

type hashable struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

func newID(inputs []Input, outputs []Output) ([]byte, error) {
	hashable := hashable{
		Inputs:  inputs,
		Outputs: outputs,
	}
	raw, err := json.Marshal(hashable)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to serialize input %#v and output %#v", inputs, outputs)
	}
	hash := sha256.Sum256(raw)
	return hash[:], nil
}

func NewTransaction(inputs []Input, outputs []Output) (*Transaction, error) {
	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create id")
	}
	return &Transaction{
		ID:      id,
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

type Input struct {
	TransactionID []byte
	Vout          int
	PublicKey     []byte
	Signature     []byte
}

type Output struct {
	Value         int
	PublicKeyHash []byte
}

type Transactions []Transaction

func (txs Transactions) Hash() []byte {
	var result []byte
	for _, tx := range txs {
		result = append(result, tx.ID...)
	}
	hash := sha256.Sum256(result)
	return hash[:]
}
