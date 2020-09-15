package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type CastVote func(from, to, signature, verifier []byte) (Transaction, error)

type SaveTransaction func(Transaction) error

type DeleteTransaction func(Transaction) error

type Transaction struct {
	ID      []byte  `json:"id"`
	Inputs  Inputs  `json:"inputs"`
	Outputs Outputs `json:"outputs"`
}

var ErrInsufficientVotes = errors.New("Not enough votes available")

func (tx Transaction) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ID: %x\n", tx.ID))
	builder.WriteString("Inputs:\n")
	for _, in := range tx.Inputs {
		builder.WriteString(fmt.Sprintf("\tFrom: %x\n", in.PublicKeyHash))
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
	signature, err := wallet.Sign(signable, creator.PrivateKey)
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
			Vout:          -1,
			PublicKeyHash: creator.PublicKeyHash(),
			Signature:     signature,
			Verifier:      creator.PublicKey,
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

func (t Transaction) UTXOs() (utxos []UTXO) {
	for i, out := range t.Outputs {
		utxos = append(utxos, UTXO{
			PublicKeyHash: out.PublicKeyHash,
			TransactionID: t.ID,
			Value:         out.Value,
			Vout:          i,
		})
	}
	return
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

func (txs Transactions) FindTransactionTo(publicKeyHash []byte) (Transaction, bool) {
	for _, tx := range txs {
		for _, output := range tx.Outputs {
			if bytes.Compare(output.PublicKeyHash, publicKeyHash) == 0 {
				return tx, true
			}
		}
	}
	return Transaction{}, false
}

func (txs Transactions) Find(criteria func(Transaction) bool) (Transaction, bool) {
	for _, tx := range txs {
		if criteria(tx) {
			return tx, true
		}
	}
	return Transaction{}, false
}
