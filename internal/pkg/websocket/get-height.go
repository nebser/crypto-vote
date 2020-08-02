package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type GetHeightFn func() (int, error)

type response struct {
	Result getHeightResult `json:"result"`
	Error  Error           `json:"error"`
}

type getHeightResult struct {
	Height int `json:"height"`
}

func GetHeight(conn *websocket.Conn) GetHeightFn {
	return func() (int, error) {
		payload := Command{
			Type: GetBlockchainHeightCommand,
		}
		if err := conn.WriteJSON(payload); err != nil {
			return 0, errors.Wrapf(err, "Failed to marshal payload into json %#v", payload)
		}
		var r response
		if err := conn.ReadJSON(&r); err != nil {
			return 0, errors.Wrap(err, "Failed to unmarshal response into get height response")
		}
		return r.Result.Height, nil
	}
}
