package operations

import (
	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type GetMissingBlocksFn func(lastBlock []byte) ([][]byte, error)

type getMissingBlocksPayload struct {
	LastBlock []byte `json:"lastBlock"`
}

type getMissingBlocksResult struct {
	Blocks [][]byte `json:"blocks"`
}

func GetMissingBlocks(conn *websocket.Conn) GetMissingBlocksFn {
	return func(lastBlock []byte) ([][]byte, error) {
		payload := operation{
			Message: _websocket.GetMissingBlocksMessage,
			Body:    getMissingBlocksPayload{LastBlock: lastBlock},
		}
		var r getMissingBlocksResult
		if err := call(conn, payload, &r); err != nil {
			return nil, err
		}
		return r.Blocks, nil
	}
}
