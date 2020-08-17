package handlers

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type getMissingBlocksResponse struct {
	Blocks [][]byte `json:"blocks"`
}

type getMissingBlocksPayload struct {
	LastBlock []byte `json:"lastBlock"`
}

func getMissingBlocks(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn, current []byte, final []byte) ([][]byte, error) {
	if len(current) == 0 || bytes.Compare(current, final) == 0 {
		return nil, nil
	}
	block, err := getBlock(current)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve block %x", current)
	}
	blocks, err := getMissingBlocks(getTip, getBlock, block.Header.Prev, final)
	if err != nil {
		return nil, err
	}
	return append(blocks, current), nil
}

func GetMissingBlocks(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn) websocket.Handler {
	return func(ping websocket.Ping) (*websocket.Pong, error) {
		var payload getMissingBlocksPayload
		if err := json.Unmarshal(ping.Body, &payload); err != nil {
			return websocket.NewErrorPong(websocket.NewInvalidDataError(websocket.GetMissingBlocksMessage.String())), nil
		}
		result, err := getMissingBlocks(getTip, getBlock, getTip(), payload.LastBlock)
		if err != nil {
			return nil, err
		}
		log.Printf("Num of blocks %d", len(result))
		return websocket.NewResponsePong(
			getMissingBlocksResponse{
				Blocks: result,
			},
		), nil
	}
}
