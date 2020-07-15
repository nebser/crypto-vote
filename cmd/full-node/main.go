package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
)

func main() {
	db, err := bolt.Open("db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	blockchain, err := blockchain.NewBlockchain(db)()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Blockchain %#v", blockchain)
}
