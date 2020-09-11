package transaction

type Input struct {
	TransactionID []byte
	Vout          int
	PublicKeyHash []byte
	Verifier      []byte
	Signature     []byte
}

type Inputs []Input
