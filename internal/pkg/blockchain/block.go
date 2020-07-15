package blockchain

type Metadata struct {
	MagicNumber int
	Size        int
}

type Header struct {
	Version         int
	Prev            []byte
	TransactionHash []byte
	Timestamp       int64
}

type Body struct {
	TransactionsCount int
	Transactions      []string
}

type Block struct {
	Metadata Metadata
	Header   Header
	Body     Body
}
