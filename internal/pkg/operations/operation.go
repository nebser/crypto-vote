package operations

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type operation struct {
	Message   websocket.Message `json:"message"`
	Body      interface{}       `json:"body"`
	Sender    string            `json:"sender,omitempty"`
	Signature string            `json:"signature,omitempty"`
}

type signable struct {
	Body     interface{}       `json:"body"`
	Sender   string            `json:"sender,omitempty"`
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
