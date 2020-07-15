package main

import (
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/repository"
)

func main() {
	repository, err := repository.New("db")
	if err != nil {
		log.Fatal(err)
	}
	blockchain, err := blockchain.NewBlockchain(repository.GetTip, repository.InitBlockchain)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Blockchain %#v", blockchain)
}
