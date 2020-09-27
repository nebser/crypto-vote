package transaction

type Input struct {
	TransactionID []byte
	Vout          int
	PublicKeyHash []byte
	Verifier      []byte
	Signature     []byte
}

type Inputs []Input

func (ins Inputs) Find(criteria func(Input) bool) (Input, bool) {
	for _, in := range ins {
		if criteria(in) {
			return in, true
		}
	}
	return Input{}, false
}
