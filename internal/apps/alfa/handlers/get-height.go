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
	return func(_ websocket.Request) (websocket.Response, error) {
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return websocket.Response{}, errors.Wrap(err, "Failed to get height")
		}
		return websocket.Response{
			Result: getHeightResponse{Height: height},
		}, nil
	}
}
