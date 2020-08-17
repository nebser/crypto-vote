package handlers

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type getBlockPayload struct {
	Hash []byte `json:"hash"`
}

type getBlockResponse struct {
	Block blockchain.Block `json:"block"`
}

func GetBlock(getBlock blockchain.GetBlockFn) websocket.Handler {
	return func(ping websocket.Ping) (*websocket.Pong, error) {
		var p getBlockPayload
		if err := json.Unmarshal(ping.Body, &p); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal data %s into payload", ping.Body)
		}
		block, err := getBlock(p.Hash)
		switch {
		case err != nil:
			return nil, errors.Wrapf(err, "Failed to retrieve block %s", p.Hash)
		case block == nil:
			return websocket.NewErrorPong(websocket.NewBlockNotFoundError(p.Hash)), nil
		default:
			return websocket.NewResponsePong(
				getBlockResponse{
					Block: *block,
				},
			), nil
		}
	}
}
