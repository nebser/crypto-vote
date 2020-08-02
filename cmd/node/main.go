package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

type response struct {
	Result Result            `json:"result"`
	Error  *_websocket.Error `json:"error"`
}

type Result struct {
	Height int `json:"height"`
}

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:10000",
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}
	defer conn.Close()

	payload := _websocket.Command{
		Type: _websocket.GetBlockchainHeightCommand,
	}
	if err := conn.WriteJSON(payload); err != nil {
		log.Fatal(err)
	}
	var r response
	if err := conn.ReadJSON(&r); err != nil {
		log.Fatal(err)
	}
	log.Printf("Received %#v\n", r)
}
