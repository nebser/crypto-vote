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
	return func(payload json.RawMessage) (websocket.Response, error) {
		var p getBlockPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return websocket.Response{}, errors.Wrapf(err, "Failed to unmarshal data %s into payload", payload)
		}
		block, err := getBlock(p.Hash)
		switch {
		case err != nil:
			return websocket.Response{}, errors.Wrapf(err, "Failed to retrieve block %s", p.Hash)
		case block == nil:
			return websocket.Response{
				Error: websocket.NewBlockNotFoundError(p.Hash),
			}, nil
		default:
			return websocket.Response{
				Result: getBlockResponse{
					Block: *block,
				},
			}, nil
		}
	}
}
