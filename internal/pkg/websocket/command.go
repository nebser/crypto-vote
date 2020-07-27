package websocket

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	GetBlockchainHeightCommand CommandType = "get-blockchain-height"
	CloseConnectionCommand     CommandType = "close-connection"
)

type CommandType string

func (c *CommandType) UnmarshalJSON(b []byte) error {
	var help string
	if err := json.Unmarshal(b, &help); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal %s into command type", b)
	}
	command := CommandType(help)
	switch command {
	case GetBlockchainHeightCommand:
		*c = command
		return nil
	default:
		return errors.Errorf("Invalid value specified for command type %s", command)
	}
}

type Command struct {
	Type      CommandType     `json:"type"`
	Body      json.RawMessage `json:"body"`
	Signature string          `json:"signature"`
	Sender    string          `json:"sender"`
}
