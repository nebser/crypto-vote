package handlers

import (
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type getHeightResponse struct {
	Height int `json:"height"`
}

func GetHeightHandler(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn) websocket.Handler {
	return func(websocket.Ping, string) (*websocket.Pong, error) {
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get height")
		}
		return websocket.NewResponsePong(
			getHeightResponse{Height: height},
		), nil
	}
}
