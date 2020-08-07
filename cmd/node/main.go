package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/nebser/crypto-vote/internal/pkg/operations"
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

	height, err := operations.GetHeight(conn)()
	if err != nil {
		log.Fatalf("Fatal error occurred %s\n", err)
	}
	log.Printf("Received height %d\n", height)
	blocks, err := operations.GetMissingBlocks(conn)(nil)
	if err != nil {
		log.Fatalf("Fatal error occurred %s\n", err)
	}
	log.Println("Received blocks")
	getBlock := operations.GetBlock(conn)
	for _, b := range blocks {
		block, err := getBlock(b)
		if err != nil {
			log.Printf("Error occurred %s\n", err)
		} else {
			log.Printf("Block found %s", block)
		}
	}
	nodes, err := operations.Register(conn)("1")
	if err != nil {
		log.Fatalf("Failed to get nodes %s\n", err)
	}
	log.Printf("Nodes %#v\n", nodes)
}
