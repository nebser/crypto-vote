package transaction

type Input struct {
	TransactionID []byte
	Vout          int
	PublicKeyHash []byte
	Signature     []byte
}

type Inputs []Input
