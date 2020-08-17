package handlers

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type registerPayload struct {
	NodeID string `json:"nodeId"`
}

type registerResponse struct {
	Nodes blockchain.Nodes `json:"nodes"`
}

func Register(saveNode blockchain.SaveNodeFn, getNodes blockchain.GetNodesFn) websocket.Handler {
	return func(ping websocket.Ping) (*websocket.Pong, error) {
		var p registerPayload
		if err := json.Unmarshal(ping.Body, &p); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal data %s into payload", ping.Body)
		}
		node := blockchain.Node{ID: p.NodeID, Type: blockchain.RegularNodeType}
		if err := saveNode(node); err != nil {
			return nil, errors.Wrapf(err, "Failed to save node %#v", node)
		}
		nodes, err := getNodes()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get nodes")
		}
		return websocket.NewResponsePong(
			registerResponse{
				Nodes: nodes,
			},
		), nil
	}
}
