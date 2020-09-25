package transaction

import (
	"bytes"
	"crypto/sha256"
	"strings"
)

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

func (txs Transactions) Len() int {
	return len(txs)
}

func (txs Transactions) Less(i, j int) bool {
	return txs[i].Timestamp < txs[j].Timestamp
}

func (txs Transactions) Swap(i, j int) {
	txs[i], txs[j] = txs[j], txs[i]
}
