package repository

import (
	"encoding/base64"
	"encoding/json"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/pkg/errors"
)

func transactionsBucket() []byte {
	return []byte("ether")
}

type tx struct {
	ID        string              `json:"id"`
	Inputs    []transactionInput  `json:"inputs"`
	Outputs   []transactionOutput `json:"outputs"`
	Timestamp int64               `json:"timestamp"`
}

func (t tx) toTransaction() transaction.Transaction {
	inputs := transaction.Inputs{}
	for _, in := range t.Inputs {
		inputs = append(inputs, in.toInput())
	}
	outputs := transaction.Outputs{}
	for _, out := range t.Outputs {
		outputs = append(outputs, out.toOutput())
	}
	id, _ := base64.StdEncoding.DecodeString(t.ID)
	return transaction.Transaction{
		ID:        id,
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: t.Timestamp,
	}
}

func newTX(transaction transaction.Transaction) tx {
	inputs := []transactionInput{}
	for _, in := range transaction.Inputs {
		inputs = append(inputs, newTransactionInput(in))
	}
	outputs := []transactionOutput{}
	for _, out := range transaction.Outputs {
		outputs = append(outputs, newTransactionOutput(out))
	}
	return tx{
		ID:        base64.StdEncoding.EncodeToString(transaction.ID),
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: transaction.Timestamp,
	}
}

type transactionInput struct {
	TransactionID string `json:"transactionId"`
	Vout          int    `json:"vout"`
	PublicKeyHash string `json:"publicKeyHash"`
	Signature     string `json:"signature"`
	Verifier      string `json:"verifier"`
}

func (ti transactionInput) toInput() transaction.Input {
	transactionID, _ := base64.StdEncoding.DecodeString(ti.TransactionID)
	publicKeyHash, _ := base64.StdEncoding.DecodeString(ti.PublicKeyHash)
	signature, _ := base64.StdEncoding.DecodeString(ti.Signature)
	verifier, _ := base64.StdEncoding.DecodeString(ti.Verifier)
	return transaction.Input{
		TransactionID: transactionID,
		Vout:          ti.Vout,
		PublicKeyHash: publicKeyHash,
		Signature:     signature,
		Verifier:      verifier,
	}
}

func newTransactionInput(input transaction.Input) transactionInput {
	return transactionInput{
		TransactionID: base64.StdEncoding.EncodeToString(input.TransactionID),
		Vout:          input.Vout,
		PublicKeyHash: base64.StdEncoding.EncodeToString(input.PublicKeyHash),
		Signature:     base64.StdEncoding.EncodeToString(input.Signature),
		Verifier:      base64.StdEncoding.EncodeToString(input.Verifier),
	}
}

type transactionOutput struct {
	Value         int    `json:"value"`
	PublicKeyHash string `json:"publicKeyHash"`
}

func (to transactionOutput) toOutput() transaction.Output {
	publicKeyHash, _ := base64.StdEncoding.DecodeString(to.PublicKeyHash)
	return transaction.Output{
		Value:         to.Value,
		PublicKeyHash: publicKeyHash,
	}
}

func newTransactionOutput(output transaction.Output) transactionOutput {
	return transactionOutput{
		Value:         output.Value,
		PublicKeyHash: base64.StdEncoding.EncodeToString(output.PublicKeyHash),
	}
}

func CastVote(db *bolt.DB) transaction.CastVote {
	return func(from, to, signature, verifier []byte) (transaction.Transaction, error) {
		var result transaction.Transaction
		err := db.Update(func(tx *bolt.Tx) error {
			utxos, err := getUTXOsByPublicKey(tx, from)
			switch {
			case err != nil:
				return errors.Wrapf(err, "Failed to retrieve utxos for %x", from)
			case len(utxos) == 0:
				return transaction.ErrInsufficientVotes
			}
			usedUTXO := utxos[0]
			inputs := transaction.Inputs{
				{
					PublicKeyHash: from,
					Signature:     signature,
					TransactionID: usedUTXO.TransactionID,
					Vout:          usedUTXO.Vout,
					Verifier:      verifier,
				},
			}
			outputs := transaction.Outputs{
				transaction.Output{
					PublicKeyHash: to,
					Value:         transaction.VoteValue,
				},
			}
			if usedUTXO.Value > transaction.VoteValue {
				outputs = append(outputs, transaction.Output{
					PublicKeyHash: from,
					Value:         usedUTXO.Value - transaction.VoteValue,
				})
			}
			tr, err := transaction.NewTransaction(inputs, outputs)
			if err != nil {
				return errors.Wrap(err, "Failed to create new transaction")
			}
			if err := saveTransaction(tx, *tr); err != nil {
				return errors.Wrap(err, "Failed to save transaction")
			}
			result = *tr
			return nil
		})
		return result, err
	}
}

func saveTransaction(tx *bolt.Tx, transaction transaction.Transaction) error {
	b := tx.Bucket(transactionsBucket())
	if b == nil {
		created, err := tx.CreateBucket(transactionsBucket())
		if err != nil {
			return errors.Wrapf(err, "Failed to create bucket %s", transactionsBucket())
		}
		b = created
	}
	raw, err := json.Marshal(newTX(transaction))
	if err != nil {
		return errors.Wrapf(err, "Failed to serialize transaction %#v", transactionsBucket())
	}
	if err := b.Put(transaction.ID, raw); err != nil {
		return errors.Wrapf(err, "Failed to save transaction %s", transaction)
	}
	return nil
}

func getInputSum(tx *bolt.Tx, tr transaction.Transaction) (int, error) {
	sum := 0
	for _, in := range tr.Inputs {
		utxo, err := getTransactionUTXO(tx, in.TransactionID, in.Vout)
		switch {
		case err != nil:
			return 0, errors.Wrapf(err, "Failed to get transaction utxo %x %d", in.TransactionID, in.Vout)
		case utxo == nil:
			return 0, transaction.ErrUTXONotFound
		}
		sum += utxo.Value
	}
	return sum, nil
}

func SaveTransaction(db *bolt.DB) transaction.SaveTransaction {
	return func(tr transaction.Transaction) error {
		return db.Update(func(tx *bolt.Tx) error {
			sum, err := getInputSum(tx, tr)
			if err != nil {
				return err
			}
			if sum != tr.Outputs.Sum() {
				return errors.Errorf("Sums of inputs (%d) and outputs (%d) are not the same", sum, tr.Outputs.Sum())
			}
			if err := saveTransaction(tx, tr); err != nil {
				return errors.Wrap(err, "Failed to save transaction")
			}
			return nil
		})
	}
}

func GetTransactions(db *bolt.DB) transaction.GetTransactionsFn {
	return func() (transaction.Transactions, error) {
		var transactions transaction.Transactions
		err := db.View(func(tx_ *bolt.Tx) error {
			b := tx_.Bucket(transactionsBucket())
			if b == nil {
				return nil
			}
			cursor := b.Cursor()
			for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
				var t tx
				if err := json.Unmarshal(value, &t); err != nil {
					return errors.Wrapf(err, "Failed to unmarshal transaction %s", value)
				}
				transactions = append(transactions, t.toTransaction())
			}
			sort.Sort(transactions)
			return nil

		})
		return transactions, err
	}
}

func deleteTransaction(tx *bolt.Tx, transaction transaction.Transaction) error {
	b := tx.Bucket(transactionsBucket())
	if b == nil {
		return nil
	}
	if err := b.Delete(transaction.ID); err != nil {
		return errors.Wrapf(err, "Failed to delete transaction %s", transaction)
	}
	return nil
}

func deleteTransactions(tx *bolt.Tx, transactions transaction.Transactions) error {
	for _, transaction := range transactions {
		if err := deleteTransaction(tx, transaction); err != nil {
			return errors.Wrap(err, "Failed to delete transactions")
		}
	}
	return nil
}
