package operations

import (
	"github.com/gorilla/websocket"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type RegisterFn func(nodeID string) (blockchain.Nodes, error)

type registerPayload struct {
	NodeID string `json:"nodeId"`
}

type registerResponse struct {
	Result registerResult    `json:"result"`
	Error  *_websocket.Error `json:"error"`
}

type registerResult struct {
	Nodes blockchain.Nodes `json:"nodes"`
}

func Register(conn *websocket.Conn) RegisterFn {
	return func(nodeID string) (blockchain.Nodes, error) {
		payload := operation{
			Type: _websocket.RegisterCommand,
			Body: registerPayload{
				NodeID: nodeID,
			},
		}
		var r registerResponse
		if err := send(conn, payload, &r); err != nil {
			return nil, errors.Wrapf(err, "Failed to send operation %#v", payload)
		}
		if r.Error != nil {
			return nil, errors.Errorf("Failed to register to node list %s", r.Error)
		}
		return r.Result.Nodes, nil
	}
}
