package handlers

import (
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func ForgeBlock(
	getTip blockchain.GetTipFn,
	getBlock blockchain.GetBlockFn,
	forgeBlock blockchain.ForgeBlockFn,
	getTransactions transaction.GetTransactionsFn,
	newStakeTransaction transaction.NewStakeTransactionFn,
	broadcast websocket.BroadcastFn,
) websocket.Handler {
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
			return nil, errors.Errorf("Cannot forge block because blockchain height is not high enough(%d)", height)
		}
		stake, err := newStakeTransaction()
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to create stake transaction")
		}
		log.Printf("Stake transaction %s", *stake)
		transactions, err := getTransactions()
		switch {
		case err != nil:
			return nil, errors.Wrap(err, "Failed to retrieve transactions")
		case len(transactions) == 0:
			log.Println("No transactions to use for forging")
			return websocket.NewNoActionPong(), nil
		}
		block, err := forgeBlock(append(transaction.Transactions{*stake}, transactions...))
		switch {
		case err != nil:
			return nil, errors.Wrap(err, "Failed to forge block")
		case block == nil:
			log.Printf("Block is not forged because there are no transactions")
			return websocket.NewNoActionPong(), nil
		}
		log.Println("Forged block")
		newBlock, err := getBlock(getTip())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to return block")
		}
		broadcast(websocket.Pong{
			Message: websocket.BlockForgedMessage,
			Body: websocket.BlockForgedBody{
				Height: height + 1,
				Block:  *newBlock,
			},
		})
		log.Println("Sent forged block")
		return websocket.NewNoActionPong(), nil
	}
}
