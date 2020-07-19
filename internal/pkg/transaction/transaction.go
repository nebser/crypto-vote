package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type Transaction struct {
	ID      []byte  `json:"id"`
	Inputs  Inputs  `json:"inputs"`
	Outputs Outputs `json:"outputs"`
}

func (tx Transaction) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ID: %x\n", tx.ID))
	builder.WriteString("Inputs:\n")
	for _, in := range tx.Inputs {
		builder.WriteString(fmt.Sprintf("\tFrom: %x\n", in.PublicKey))
		builder.WriteString(fmt.Sprintf("\tSignature: %x\n", in.Signature))
	}
	builder.WriteString("Outputs:\n")
	for _, out := range tx.Outputs {
		builder.WriteString(fmt.Sprintf("\tTo: %x\n", out.PublicKeyHash))
		builder.WriteString(fmt.Sprintf("\tValue: %d\n", out.Value))
	}
	return builder.String()
}

type hashable struct {
	Inputs  Inputs  `json:"inputs"`
	Outputs Outputs `json:"outputs"`
}

func newID(inputs Inputs, outputs Outputs) ([]byte, error) {
	hashable := hashable{
		Inputs:  inputs,
		Outputs: outputs,
	}
	return hash(hashable)
}

func hash(data interface{}) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to serialize data %#v", data)
	}
	hash := sha256.Sum256(raw)
	return hash[:], nil
}

func NewTransaction(inputs Inputs, outputs Outputs) (*Transaction, error) {
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

func NewBaseTransaction(creator wallet.Wallet, recipientAddress string) (*Transaction, error) {
	recipientKeyHash := wallet.ExtractPublicKeyHash(recipientAddress)
	signable := signable{
		Recipient: recipientKeyHash,
		Sender:    creator.PublicKeyHash(),
		Value:     1,
	}
	signature, err := sign(signable, creator.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to sign base transaction")
	}
	outputs := Outputs{
		{
			Value:         1,
			PublicKeyHash: recipientKeyHash,
		},
	}
	inputs := Inputs{
		{
			Vout:      -1,
			PublicKey: creator.PublicKey,
			Signature: signature,
		},
	}
	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create transaction id")
	}
	return &Transaction{
		ID:      id,
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func (t Transaction) IsBase() bool {
	return len(t.Inputs) == 1 && len(t.Outputs) == 1 && t.Inputs[0].Vout == -1
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

func (txs Transactions) String() string {
	builder := strings.Builder{}
	builder.WriteString("-----START TRANSACTIONS-----\n")
	for _, tx := range txs {
		builder.WriteString(tx.String())
	}
	builder.WriteString("-----END TRANSACTIONS-----\n")
	return builder.String()
}
