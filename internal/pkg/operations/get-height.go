package operations

import (
	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type GetHeightFn func() (int, error)

type response struct {
	Result getHeightResult   `json:"result"`
	Error  *_websocket.Error `json:"error"`
}

type getHeightResult struct {
	Height int `json:"height"`
}

func GetHeight(conn *websocket.Conn) GetHeightFn {
	return func() (int, error) {
		payload := operation{
			Type: _websocket.GetBlockchainHeightCommand,
		}
		var r response
		if err := send(conn, payload, &r); err != nil {
			return 0, err
		}
		if r.Error != nil {
			return 0, errors.Errorf("Failed to get height %s", r.Error.Message)
		}
		return r.Result.Height, nil
	}
}
