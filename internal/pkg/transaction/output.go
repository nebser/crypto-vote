package transaction

type Output struct {
	Value         int
	PublicKeyHash []byte
}

type Outputs []Output
