package operations

import (
	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type GetMissingBlocksFn func(lastBlock []byte) ([][]byte, error)

type getMissingBlocksPayload struct {
	LastBlock []byte `json:"lastBlock"`
}

type getMissingBlocksResult struct {
	Blocks [][]byte `json:"blocks"`
}

type getMissingBlocksResponse struct {
	Result getMissingBlocksResult `json:"result"`
	Error  *_websocket.Error      `json:"error"`
}

func GetMissingBlocks(conn *websocket.Conn) GetMissingBlocksFn {
	return func(lastBlock []byte) ([][]byte, error) {
		payload := operation{
			Type: _websocket.GetMissingBlocksCommand,
			Body: getMissingBlocksPayload{LastBlock: lastBlock},
		}
		var r getMissingBlocksResponse
		if err := send(conn, payload, &r); err != nil {
			return nil, errors.Wrapf(err, "Failed to send operation %#v", payload)
		}
		if r.Error != nil {
			return nil, errors.Errorf("Failed to get missing blocks %s", r.Error.Message)
		}
		return r.Result.Blocks, nil
	}
}
