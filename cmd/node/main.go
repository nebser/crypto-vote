package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

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

	height, err := _websocket.GetHeight(conn)()
	if err != nil {
		log.Fatalf("Fatal error occurred %s\n", err)
	}
	log.Printf("Received height %d\n", height)
	blocks, err := _websocket.GetMissingBlocks(conn)(nil)
	if err != nil {
		log.Fatalf("Fatal error occurred %s\n", err)
	}
	log.Println("Received blocks")
	for _, b := range blocks {
		log.Printf("Block %x", b)
	}
}
