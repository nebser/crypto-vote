package operations

import (
	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type GetHeightFn func() (int, error)

type getHeightResult struct {
	Height int `json:"height"`
}

func GetHeight(conn *websocket.Conn) GetHeightFn {
	return func() (int, error) {
		payload := operation{
			Type: _websocket.GetBlockchainHeightCommand,
		}
		var r getHeightResult
		if err := call(conn, payload, &r); err != nil {
			return 0, err
		}
		return r.Height, nil
	}
}
