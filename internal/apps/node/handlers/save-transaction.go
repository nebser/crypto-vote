package handlers

import (
	"encoding/json"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func SaveTransaction(save transaction.SaveTransaction, verifier wallet.VerifierFn) websocket.Handler {
	return func(ping websocket.Ping, _ string) (*websocket.Pong, error) {
		log.Println("STARTED SAVING")
		var p websocket.SaveTransactionBody
		if err := json.Unmarshal(ping.Body, &p); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal data %s into payload", ping.Body)
		}
		switch ok, err := verifier(ping, ping.Signature, ping.Sender); {
		case err != nil:
			return nil, errors.Wrap(err, "Failed to verify transaction")
		case !ok:
			return websocket.NewErrorPong(websocket.NewInvalidTransactionError()), nil
		}
		log.Println("TRANSACTION VERIFIED")
		if err := save(p.Transaction); err != nil {
			return nil, errors.Wrapf(err, "Failed to save transaction %s", p.Transaction)
		}
		log.Println("SAVED TRANSACTION")
		return websocket.NewNoActionPong(), nil
	}
}
