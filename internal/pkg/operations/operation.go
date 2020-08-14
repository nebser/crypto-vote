package operations

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type operation struct {
	Type      websocket.CommandType `json:"type"`
	Body      interface{}           `json:"body"`
	Sender    string                `json:"sender"`
	Signature string                `json:"signature"`
}

type signable struct {
	Body   interface{} `json:"body"`
	Sender string      `json:"sender"`
}

func (op operation) Signable() ([]byte, error) {
	s := signable{
		Body:   op.Body,
		Sender: op.Sender,
	}
	return json.Marshal(s)
}
