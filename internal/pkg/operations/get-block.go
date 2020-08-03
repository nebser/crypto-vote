package operations

import (
	"github.com/gorilla/websocket"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type GetBlockFn func(blockHash []byte) (blockchain.Block, error)

type getBlockPayload struct {
	Hash []byte `json:"hash"`
}

type getBlockResult struct {
	Block blockchain.Block `json:"block"`
}

type getBlockResponse struct {
	Result getBlockResult    `json:"result"`
	Error  *_websocket.Error `json:"error"`
}

type ErrBlockNotFound string

func (e ErrBlockNotFound) Error() string {
	return string(e)
}

func GetBlock(conn *websocket.Conn) GetBlockFn {
	return func(blockHash []byte) (blockchain.Block, error) {
		payload := operation{
			Type: _websocket.GetBlockCommand,
			Body: getBlockPayload{Hash: blockHash},
		}
		var r getBlockResponse
		if err := send(conn, payload, &r); err != nil {
			return blockchain.Block{}, err
		}
		if r.Error != nil {
			if *r.Error == *_websocket.NewBlockNotFoundError(blockHash) {
				return blockchain.Block{}, ErrBlockNotFound(r.Error.Message)
			}
			return blockchain.Block{}, errors.Errorf("Failed to get block %s", blockHash)
		}
		return r.Result.Block, nil
	}
}
