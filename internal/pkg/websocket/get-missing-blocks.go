package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
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
	Error  Error                  `json:"error"`
}

func GetMissingBlocks(conn *websocket.Conn) GetMissingBlocksFn {
	return func(lastBlock []byte) ([][]byte, error) {
		raw, err := json.Marshal(getMissingBlocksPayload{LastBlock: lastBlock})
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to marshal last block %x", lastBlock)
		}
		command := Command{
			Type: GetMissingBlocksCommand,
			Body: raw,
		}
		if err := conn.WriteJSON(command); err != nil {
			return nil, errors.Wrapf(err, "Failed to marshal payload into json %#v", command)
		}
		var r getMissingBlocksResponse
		if err := conn.ReadJSON(&r); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal response into get missing blocks response")
		}
		return r.Result.Blocks, nil
	}
}
