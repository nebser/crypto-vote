package handlers

import (
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func ForgeBlock(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn, forgeBlock blockchain.ForgeBlockFn, getTransactions transaction.GetTransactionsFn) websocket.Handler {
	return func(ping websocket.Ping, _ string) (*websocket.Pong, error) {
		var body websocket.ForgeBlockBody
		if err := json.Unmarshal(ping.Body, &body); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal forge block message body %s", ping.Body)
		}
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to retrieve block height")
		}
		if height < body.Height {
			return nil, errors.Errorf("Cannot forge block because blockchain height (%s)")
		}
		transactions, err := getTransactions()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to retrieve transactions")
		}
		block, err := forgeBlock(transactions)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to forge block %s", err)
		}
		log.Printf("Forged block %s\n", *block)
		return websocket.NewNoActionPong(), nil
	}
}
