package transaction

import (
	"encoding/json"
)

type signable struct {
	Sender    []byte `json:"sender"`
	Recipient []byte `json:"recipient"`
	Value     int    `json:"value"`
}

func (s signable) Signable() ([]byte, error) {
	return json.Marshal(s)
}
