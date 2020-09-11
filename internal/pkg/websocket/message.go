package websocket

import (
	"encoding/json"
	"fmt"
)

type Message int

const (
	GetBlockchainHeightMessage Message = iota + 1
	CloseConnectionMessage
	GetMissingBlocksMessage
	GetBlockMessage
	RegisterMessage
	ErrorMessage
	ResponseMessage
	TransactionReceivedMessage
)

func (m Message) String() string {
	switch m {
	case GetBlockchainHeightMessage:
		return "get-blockchain-height"
	case CloseConnectionMessage:
		return "close-connection"
	case GetMissingBlocksMessage:
		return "get-missing-blocks"
	case GetBlockMessage:
		return "get-block"
	case RegisterMessage:
		return "register"
	case ErrorMessage:
		return "error"
	case ResponseMessage:
		return "response"
	case TransactionReceivedMessage:
		return "transaction-received"
	default:
		return fmt.Sprintf("Unknown message %d", m)
	}
}

type Ping struct {
	Message   Message         `json:"message"`
	Body      json.RawMessage `json:"body"`
	Signature string          `json:"signature"`
	Sender    string          `json:"sender"`
}

type signablePing struct {
	Body    json.RawMessage `json:"body"`
	Sender  string          `json:"sender"`
	Message Message         `json:"message"`
}

func (p Ping) Signable() ([]byte, error) {
	s := signablePing{
		Body:    p.Body,
		Message: p.Message,
		Sender:  p.Sender,
	}
	return json.Marshal(s)
}

type Pong struct {
	Message   Message     `json:"message"`
	Body      interface{} `json:"body"`
	Signature string      `json:"signature"`
	Sender    string      `json:"sender"`
}

func NewErrorPong(e Error) *Pong {
	return &Pong{
		Message: ErrorMessage,
		Body:    e,
	}
}

func NewResponsePong(body interface{}) *Pong {
	return &Pong{
		Message: ResponseMessage,
		Body:    body,
	}
}
