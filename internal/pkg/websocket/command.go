package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	GetBlockchainHeightCommand CommandType = iota + 1
	CloseConnectionCommand
	GetMissingBlocksCommand
	GetBlockCommand
)

type CommandType int

func (c CommandType) String() string {
	switch c {
	case GetBlockchainHeightCommand:
		return "get-blockchain-height"
	case CloseConnectionCommand:
		return "close-connection"
	case GetMissingBlocksCommand:
		return "get-missing-blocks"
	case GetBlockCommand:
		return "get-block"
	default:
		return fmt.Sprintf("Unknown command %d", c)
	}
}

func (c *CommandType) UnmarshalJSON(b []byte) error {
	var help int
	if err := json.Unmarshal(b, &help); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal %s into command type", b)
	}
	command := CommandType(help)
	switch command {
	case GetBlockchainHeightCommand, CloseConnectionCommand, GetMissingBlocksCommand, GetBlockCommand:
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
