package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type CastVote func(from, to, signature, verifier []byte) (Transaction, error)

type SaveTransaction func(Transaction) error

type GetTransactionsFn func() (Transactions, error)

type DeleteTransaction func(Transaction) error

type NewStakeTransactionFn func() (*Transaction, error)

type VerifyTransctionFn func(Transaction) bool

type IsStakeTransactionFn func(Transaction) bool

type NewReturnStakeTransactionFn func(Transaction) (*Transaction, error)

const VoteValue = 10

type Transaction struct {
	ID        []byte  `json:"id"`
	Inputs    Inputs  `json:"inputs"`
	Outputs   Outputs `json:"outputs"`
	Timestamp int64   `json:"timestamp"`
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
	Inputs    Inputs  `json:"inputs"`
	Outputs   Outputs `json:"outputs"`
	Timestamp int64   `json:"timestamp"`
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
		ID:        id,
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now().Unix(),
	}, nil
}

func NewStakeTransaction(getUTXOs GetUTXOsByPublicKeyFn, signer wallet.Signer, stakeCreator wallet.Wallet, stakeholder []byte) NewStakeTransactionFn {
	return func() (*Transaction, error) {
		utxos, err := getUTXOs(stakeCreator.PublicKeyHash())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to retrieve utxos for stake tx for %x", stakeCreator.PublicKeyHash())
		}
		target := utxos.Sum() / 2
		if target < VoteValue/2 {
			return nil, ErrCantForge
		}
		sum := 0
		var inputs Inputs
		for _, utxo := range utxos {
			sum += utxo.Value
			signable := signable{
				Recipient: stakeholder,
				Sender:    stakeCreator.PublicKeyHash(),
				Value:     utxo.Value,
			}
			signature, err := signer.SignRaw(signable)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to sign %#v", signable)
			}
			inputs = append(inputs, Input{
				PublicKeyHash: stakeCreator.PublicKeyHash(),
				Signature:     signature,
				TransactionID: utxo.TransactionID,
				Vout:          utxo.Vout,
				Verifier:      stakeCreator.PublicKey,
			})
			if sum >= target {
				break
			}
		}
		outputs := Outputs{
			{
				Value:         target,
				PublicKeyHash: stakeholder,
			},
		}
		if sum > target {
			outputs = append(outputs, Output{
				Value:         sum - target,
				PublicKeyHash: stakeCreator.PublicKeyHash(),
			})
		}
		return NewTransaction(inputs, outputs)
	}
}

func NewReturnStakeTransaction(w wallet.Wallet) NewReturnStakeTransactionFn {
	return func(transaction Transaction) (*Transaction, error) {
		pKeyHash := w.PublicKeyHash()
		index, found := transaction.Outputs.FindIndex(func(element Output) bool {
			return bytes.Compare(element.PublicKeyHash, pKeyHash) == 0
		})
		if !found {
			return nil, errors.New("Failed to find output transaction")
		}
		signable := signable{
			Recipient: transaction.Inputs[0].PublicKeyHash,
			Sender:    pKeyHash,
			Value:     transaction.Outputs[index].Value,
		}
		signature, err := wallet.Sign(signable, w.PrivateKey)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to sign return stake transaction")
		}
		inputs := Inputs{
			{
				PublicKeyHash: pKeyHash,
				Signature:     signature,
				TransactionID: transaction.ID,
				Verifier:      w.PublicKey,
				Vout:          index,
			},
		}
		outputs := Outputs{
			Output{
				Value:         transaction.Outputs[index].Value,
				PublicKeyHash: transaction.Inputs[0].PublicKeyHash,
			},
		}
		return NewTransaction(inputs, outputs)
	}
}

func NewBaseTransaction(creator wallet.Wallet, recipientAddress string) (*Transaction, error) {
	recipientKeyHash := wallet.ExtractPublicKeyHash(recipientAddress)
	signable := signable{
		Recipient: recipientKeyHash,
		Sender:    creator.PublicKeyHash(),
		Value:     VoteValue,
	}
	signature, err := wallet.Sign(signable, creator.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to sign base transaction")
	}
	outputs := Outputs{
		{
			Value:         VoteValue,
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

func (t Transaction) AreInputsFrom(pkeyHash []byte) bool {
	_, found := t.Inputs.Find(func(input Input) bool {
		return bytes.Compare(input.PublicKeyHash, pkeyHash) != 0
	})
	return !found
}

func VerifyTransactions(getTransactionUTXO GetTransactionUTXO, verifier wallet.VerifierFn) VerifyTransctionFn {
	return func(transaction Transaction) bool {
		for _, input := range transaction.Inputs {
			receiver, found := transaction.Outputs.Find(func(o Output) bool {
				return bytes.Compare(o.PublicKeyHash, input.PublicKeyHash) != 0
			})
			if !found {
				return false
			}
			utxo, err := getTransactionUTXO(input.TransactionID, input.Vout)
			if err != nil || utxo == nil {
				return false
			}
			signable := signable{
				Recipient: receiver.PublicKeyHash,
				Sender:    input.PublicKeyHash,
				Value:     utxo.Value,
			}
			signature := base64.StdEncoding.EncodeToString(input.Signature)
			pKey := base64.StdEncoding.EncodeToString(input.Verifier)
			if ok, err := verifier(signable, signature, pKey); err != nil || !ok {
				return false
			}
		}
		return true
	}
}

func IsStakeTransaction(alfaKeyHash []byte) IsStakeTransactionFn {
	return func(transaction Transaction) bool {
		if len(transaction.Outputs) > 2 {
			return false
		}
		_, found := transaction.Outputs.Find(func(o Output) bool {
			return bytes.Compare(o.PublicKeyHash, alfaKeyHash) == 0
		})
		if !found {
			return false
		}
		return true
	}
}
