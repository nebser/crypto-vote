package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nebser/crypto-vote/internal/pkg/repository"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"

	"github.com/nebser/crypto-vote/internal/apps/alfa"
	"github.com/nebser/crypto-vote/internal/apps/alfa/handlers"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
)

func main() {
	newOption := flag.Bool("new", false, "Should initialize new blockchain")
	privateKeyOption := flag.String("private", "alfa/key.pem", "Private key file path")
	publicKeyOption := flag.String("public", "alfa/key_pub.pem", "Public key file path")
	clientKeysDirOption := flag.String("clients", "clients", "Client key pair files directory")
	dbFileName := flag.String("db", "db", "File name to use for bolt database")

	flag.Parse()
	options := alfa.Options{
		PublicKeyFileName:  *publicKeyOption,
		PrivateKeyFileName: *privateKeyOption,
		ClientKeysDir:      *clientKeysDirOption,
	}
	if *newOption {
		if err := os.Remove(*dbFileName); err != nil {
			log.Fatalf("Failed to remove file %s", *dbFileName)
		}
	}
	db, err := bolt.Open(*dbFileName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	blockchain := blockchain.NewBlockchain(
		repository.GetTip(db),
		repository.InitBlockchain(db),
		repository.AddBlock(db),
		repository.GetBlock(db),
	)
	if *newOption {
		if err := alfa.Initialize(*blockchain, options); err != nil {
			log.Fatal(err)
		}
	}
	// listener, err := net.Listen("tcp", ":10000")
	// if err != nil {
	// 	log.Fatalf("Failed to start tcp server %s", err)
	// }
	blockchain.Print()
	router := websocket.Router{
		websocket.GetBlockchainHeightCommand: handlers.GetHeightHandler(*blockchain),
	}
	http.Handle("/", alfa.Connection(router))
	http.ListenAndServe(":10000", nil)
}
