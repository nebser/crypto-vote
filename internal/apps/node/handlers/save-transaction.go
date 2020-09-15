package handlers

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type saveTransactionPayload struct {
	Transaction transaction.Transaction `json:"transaction"`
}

func SaveTransaction() websocket.Handler {
	return func(ping websocket.Ping, _ string) (*websocket.Pong, error) {
		var p saveTransactionPayload
		if err := json.Unmarshal(ping.Body, &p); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal data %s into payload", ping.Body)
		}
		return websocket.NewResponsePong(nil), nil
	}
}
