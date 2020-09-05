package transaction

type UTXO struct {
	TransactionID []byte
	PublicKeyHash []byte
	Value         int
}

type SaveUTXO func(UTXO) error

type GetUTXOS func(publicKeyHash []byte) ([]UTXO, error)
