package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type blockForgedBody struct {
	Height int              `json:"height"`
	Block  blockchain.Block `json:"block"`
}

func BlockForged(getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn, verifyBlock blockchain.VerifyBlockFn, addNewBlock blockchain.AddNewBlockFn) websocket.Handler {
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
		if !verifyBlock(body.Block) || !body.Block.Body.Transactions[0].AreInputsFrom(hashedSender) {
			return websocket.NewDisconnectPong(), nil
		}
		switch err := addNewBlock(body.Block); {
		case errors.Is(err, blockchain.ErrInvalidBlock):
			return websocket.NewDisconnectPong(), nil
		case err != nil:
			return nil, errors.Wrap(err, "Failed to add new block to blockchain")
		default:
			log.Println("New block added")
			return websocket.NewNoActionPong(), nil
		}
	}
}
