package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/nebser/crypto-vote/internal/apps/alfa"
	"github.com/nebser/crypto-vote/internal/apps/node"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/repository"

	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/operations"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
)

func main() {
	nodeID := flag.Int("id", 0, "ID of the node [required]")
	newOption := flag.Bool("new", false, "Should initialize new blockchain")
	privateKeyOption := flag.String("private", "", "Private key file path [default is nodes/key_id.pem]")
	publicKeyOption := flag.String("public", "", "Private key file path [default is nodes/key_id_pub.pem]")
	flag.Parse()
	if *nodeID <= 0 {
		log.Fatal("NodeId must be provided and it must be greater than 0")
	}
	privateKey := *privateKeyOption
	if privateKey == "" {
		privateKey = fmt.Sprintf("nodes/n%d.pem", *nodeID)
	}
	publicKey := *publicKeyOption
	if publicKey == "" {
		publicKey = fmt.Sprintf("nodes/n%d_pub.pem", *nodeID)
	}
	dbFileName := fmt.Sprintf("db_%d", *nodeID)

	wallet, err := wallet.Import(keyfiles.KeyFiles{PrivateKeyFile: privateKey, PublicKeyFile: publicKey})
	if err != nil {
		log.Fatalf("Wallet could not be imported %s\n", err)
	}
	if *newOption {
		switch _, err := os.Stat(dbFileName); {
		case err == nil:
			if err := os.Remove(dbFileName); err != nil {
				log.Fatalf("Failed to remove file %s", dbFileName)
			}
		case err != nil && !os.IsNotExist(err):
			log.Fatalf("Failed to read stat for file %s", dbFileName)
		}

	}
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	getTip := repository.GetTip(db)
	getBlock := repository.GetBlock(db)
	if err := node.Initialize(
		operations.GetHeight(conn),
		operations.GetMissingBlocks(conn),
		operations.GetBlock(conn),
		getTip,
		getBlock,
		repository.AddBlocks(db),
	); err != nil {
		log.Fatalf("Failed to initialize node %s", err)
	}
	blockchain.PrintBlockchain(getTip, getBlock)
	nodes, err := operations.Register(conn, *wallet)(strconv.Itoa(*nodeID))
	if err != nil {
		log.Fatalf("Failed to register %s\n", err)
	}
	log.Printf("Nodes %#v\n", nodes)
	router := _websocket.Router{}
	http.Handle("/", alfa.Connection(router))
	http.ListenAndServe(fmt.Sprintf(":%d", 10000+*nodeID), nil)
}
