package websocket

import (
	"encoding/json"
)

type Request struct {
	Body      json.RawMessage
	Signature string
	Sender    string
}

type signable struct {
	Body   json.RawMessage `json:"body"`
	Sender string          `json:"sender"`
}

func (r Request) Signable() ([]byte, error) {
	s := signable{
		Body:   r.Body,
		Sender: r.Sender,
	}
	return json.Marshal(s)
}
