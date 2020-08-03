package operations

import "github.com/nebser/crypto-vote/internal/pkg/websocket"

type operation struct {
	Type      websocket.CommandType `json:"type"`
	Body      interface{}           `json:"body"`
	Sender    string                `json:"sender"`
	Signature string                `json:"signature"`
}
