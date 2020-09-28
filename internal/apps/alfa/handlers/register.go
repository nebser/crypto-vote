package handlers

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type registerPayload struct {
	NodeID string `json:"nodeId"`
}

type registerResponse struct {
	Nodes []string `json:"nodes"`
}

func Register(hub *websocket.Hub) websocket.Handler {
	return func(ping websocket.Ping, internalID string) (*websocket.Pong, error) {
		var p registerPayload
		if err := json.Unmarshal(ping.Body, &p); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal data %s into payload", ping.Body)
		}
		nodes := hub.RegisterAtomically(internalID, p.NodeID)
		return websocket.NewResponsePong(
			registerResponse{
				Nodes: nodes,
			},
		), nil
	}
}
