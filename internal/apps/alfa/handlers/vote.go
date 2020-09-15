package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nebser/crypto-vote/internal/pkg/api"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type voteBody struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Verifier  string `json:"verifier"`
	Signature string `json:"signature"`
}

func (v voteBody) Signable() ([]byte, error) {
	data := struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Value     int    `json:"value"`
	}{
		Sender:    v.Sender,
		Recipient: v.Recipient,
		Value:     1,
	}
	return json.Marshal(data)
}

func Vote(findBlock blockchain.FindBlockFn, castVote transaction.CastVote, broadcast websocket.BroadcastFn) api.Handler {
	return func(request api.Request) (api.Response, error) {
		var body voteBody
		if err := json.Unmarshal(request.Body, &body); err != nil {
			return api.InvalidDataErrorResponse(""), nil
		}
		rawPublicKey, err := base64.StdEncoding.DecodeString(body.Verifier)
		if err != nil {
			return api.InvalidDataErrorResponse("Invalid public key provided"), nil
		}
		rawSignature, err := base64.StdEncoding.DecodeString(body.Signature)
		if err != nil {
			return api.InvalidDataErrorResponse("Invalid signature provided"), nil
		}
		if !wallet.Verify(body, rawSignature, rawPublicKey) {
			return api.UnauthorizedErrorResponse("Signature does not match the payload"), nil
		}
		sender, err := base64.StdEncoding.DecodeString(body.Sender)
		if err != nil {
			return api.InvalidDataErrorResponse("Invalid sender provided"), nil
		}
		receiver, err := base64.StdEncoding.DecodeString(body.Recipient)
		if err != nil {
			return api.InvalidDataErrorResponse("Invalid recipient provided"), nil
		}

		criteria := func(b blockchain.Block) bool {
			if _, ok := b.Body.Transactions.FindTransactionTo(sender); ok {
				return true
			}
			return false
		}
		switch _, ok, err := findBlock(criteria); {
		case err != nil:
			return api.Response{}, errors.Errorf("Failed to find block. Error: %s", err)
		case !ok:
			return api.UnauthorizedErrorResponse(fmt.Sprintf("Recipient %s does not exist", body.Recipient)), nil
		default:
			log.Println("Authorized successfully")
		}
		tr, err := castVote(sender, receiver, rawSignature, rawPublicKey)
		switch {
		case err != nil && errors.Is(err, transaction.ErrInsufficientVotes):
			return api.UserAlreadyVoted(), nil
		case err != nil:
			log.Printf("Error occurred while voting %s", err)
			return api.Response{}, nil
		}
		log.Println("VOTED SUCCESSFULLY")
		broadcast(websocket.Pong{
			Message: websocket.TransactionReceivedMessage,
			Body:    tr,
		})
		return api.Response{
			Status: http.StatusOK,
		}, nil
	}
}
