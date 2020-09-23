package handlers

import (
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func ForgeBlock(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn, forgeBlock blockchain.ForgeBlockFn) websocket.Handler {
	return func(ping websocket.Ping, _ string) (*websocket.Pong, error) {
		log.Println("Checkpoint 1")
		var body websocket.ForgeBlockBody
		if err := json.Unmarshal(ping.Body, &body); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal forge block message body %s", ping.Body)
		}
		log.Println("Checkpoint 2")
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to retrieve block height")
		}
		log.Println("Checkpoint 3")
		if height < body.Height {
			return nil, errors.Errorf("Cannot forge block because blockchain height (%s)")
		}
		log.Println("Checkpoint 4")
		block, err := forgeBlock()
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to forge block %s", err)
		}
		log.Println("Checkpoint 5")
		log.Printf("Forged block %d\n", len(block.Body.Transactions))
		return websocket.NewNoActionPong(), nil
	}
}
