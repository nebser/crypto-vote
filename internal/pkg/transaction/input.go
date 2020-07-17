package transaction

type Input struct {
	TransactionID []byte
	Vout          int
	PublicKey     []byte
	Signature     []byte
}

type Inputs []Input
