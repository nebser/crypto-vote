package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type blockForgedBody struct {
	Height int              `json:"height"`
	Block  blockchain.Block `json:"block"`
}

func BlockForged(
	getTip blockchain.GetTipFn,
	getBlock blockchain.GetBlockFn,
	verifyBlock blockchain.VerifyBlockFn,
	addNewBlock blockchain.AddNewBlockFn,
	isStakeTransaction transaction.IsStakeTransactionFn,
	saveTransaction transaction.SaveTransaction,
	newReturnStakeTransaction transaction.NewReturnStakeTransactionFn,
	broadcast websocket.BroadcastFn,
) websocket.Handler {
	return func(ping websocket.Ping, _ string) (*websocket.Pong, error) {
		var body blockForgedBody
		if err := json.Unmarshal(ping.Body, &body); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarsha block forged body %s", ping.Body)
		}
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get height")
		}
		if height+1 < body.Height {
			return nil, errors.Errorf("Blockchain height is too low %d", height)
		}
		sender, err := base64.StdEncoding.DecodeString(ping.Sender)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to decode sender %s", ping.Sender)
		}
		hashedSender, err := wallet.HashedPublicKey(sender)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to extract hashed public key")
		}
		if len(body.Block.Body.Transactions) == 0 || !isStakeTransaction(body.Block.Body.Transactions[0]) {
			return websocket.NewErrorPong(websocket.NewInvalidDataError(websocket.BlockForgedMessage.String())), nil
		}
		stakeTx := body.Block.Body.Transactions[0]
		if !verifyBlock(body.Block, hashedSender) {
			if err := saveTransaction(stakeTx); err != nil {
				return nil, errors.Wrapf(err, "Failed to save stake transaction %s", stakeTx)
			}
			broadcast(websocket.Pong{
				Message: websocket.TransactionReceivedMessage,
				Body: websocket.SaveTransactionBody{
					Transaction: stakeTx,
				},
			})
			return websocket.NewDisconnectPong(), nil
		}
		returnStakeTx, err := newReturnStakeTransaction(stakeTx)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to create return stake transaction out of %s", stakeTx)
		}
		switch err := addNewBlock(body.Block); {
		case errors.Is(err, blockchain.ErrInvalidBlock):
			if err := saveTransaction(stakeTx); err != nil {
				return nil, errors.Wrapf(err, "Failed to save stake transaction %s", stakeTx)
			}
			broadcast(websocket.Pong{
				Message: websocket.TransactionReceivedMessage,
				Body: websocket.SaveTransactionBody{
					Transaction: stakeTx,
				},
			})
			log.Println("Block is invalid")
			return websocket.NewDisconnectPong(), nil
		case err != nil:
			return nil, errors.Wrap(err, "Failed to add new block to blockchain")
		default:
			log.Println("New block added")
			if err := saveTransaction(*returnStakeTx); err != nil {
				return nil, errors.Wrapf(err, "Failed to save stake transaction %s", stakeTx)
			}
			broadcast(websocket.Pong{
				Message: websocket.TransactionReceivedMessage,
				Body: websocket.SaveTransactionBody{
					Transaction: *returnStakeTx,
				},
			})
			return websocket.NewNoActionPong(), nil
		}
	}
}
