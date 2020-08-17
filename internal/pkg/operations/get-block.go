package operations

import (
	"github.com/gorilla/websocket"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type GetBlockFn func(blockHash []byte) (blockchain.Block, error)

type getBlockPayload struct {
	Hash []byte `json:"hash"`
}

type getBlockResult struct {
	Block blockchain.Block `json:"block"`
}

type ErrBlockNotFound string

func (e ErrBlockNotFound) Error() string {
	return string(e)
}

func GetBlock(conn *websocket.Conn) GetBlockFn {
	return func(blockHash []byte) (blockchain.Block, error) {
		payload := operation{
			Message: _websocket.GetBlockMessage,
			Body:    getBlockPayload{Hash: blockHash},
		}
		var r getBlockResult
		if err := call(conn, payload, &r); err != nil {
			return blockchain.Block{}, err
		}
		return r.Block, nil
	}
}
