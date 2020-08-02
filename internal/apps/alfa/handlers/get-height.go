package handlers

import (
	"encoding/json"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type getHeightResponse struct {
	Height int `json:"height"`
}

func GetHeightHandler(blockchain blockchain.Blockchain) websocket.Handler {
	return func(_ json.RawMessage) (websocket.Response, error) {
		height, err := blockchain.GetHeight()
		if err != nil {
			return websocket.Response{}, errors.Wrap(err, "Failed to get height")
		}
		return websocket.Response{
			Result: getHeightResponse{Height: height},
		}, nil
	}
}
