package operations

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type operation struct {
	Message   websocket.Message `json:"message"`
	Body      interface{}       `json:"body"`
	Sender    string            `json:"sender"`
	Signature string            `json:"signature"`
}

type signable struct {
	Body     interface{}       `json:"body"`
	Sender   string            `json:"sender"`
	Messsage websocket.Message `json:"message"`
}

func (op operation) Signable() ([]byte, error) {
	s := signable{
		Body:     op.Body,
		Sender:   op.Sender,
		Messsage: op.Message,
	}
	return json.Marshal(s)
}
