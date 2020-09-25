package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/nebser/crypto-vote/internal/apps/node"
	"github.com/nebser/crypto-vote/internal/apps/node/handlers"
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

	masterWallet, err := wallet.Import(keyfiles.KeyFiles{PrivateKeyFile: privateKey, PublicKeyFile: publicKey})
	if err != nil {
		log.Fatalf("Wallet could not be imported %s\n", err)
	}
	alfaPKey, err := wallet.LoadPublicKey("alfa/key_pub.pem")
	if err != nil {
		log.Fatalf("Failed to load public key %s", err)
	}
	encodedAlfaPkey := base64.StdEncoding.EncodeToString(alfaPKey)
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

	getTip := repository.GetTip(db)
	getBlock := repository.GetBlock(db)
	if err := node.Initialize(
		operations.GetHeight(conn),
		operations.GetMissingBlocks(conn),
		operations.GetBlock(conn),
		getTip,
		getBlock,
		repository.AddBlock(db),
	); err != nil {
		log.Fatalf("Failed to initialize node %s", err)
	}
	blockchain.PrintBlockchain(getTip, getBlock)
	nodes, err := operations.Register(conn, *masterWallet)(strconv.Itoa(*nodeID))
	if err != nil {
		log.Fatalf("Failed to register %s\n", err)
	}
	hub := _websocket.NewHub()
	router := _websocket.Router{
		_websocket.RegisterMessage: handlers.Register(hub).
			Authorized(
				blockchain.BlockchainAuthorizer(
					blockchain.FindBlock(
						repository.GetTip(db),
						repository.GetBlock(db),
					),
				),
			),
		_websocket.TransactionReceivedMessage: handlers.SaveTransaction(
			repository.SaveTransaction(db),
			wallet.VerifySignature,
		),
		_websocket.ForgeBlockMessage: handlers.ForgeBlock(
			repository.GetTip(db),
			repository.GetBlock(db),
			repository.ForgeBlock(db),
			repository.GetTransactions(db),
		).
			Authorized(
				_websocket.PublicKeyAuthorizer(
					encodedAlfaPkey,
					wallet.VerifySignature,
				),
			),
	}
	go _websocket.MaintainConnection(conn, router, hub, "0")
	if err := connectToNodes(nodes, *masterWallet, router, hub); err != nil {
		log.Fatalf("Failed to connect to nodes %s", err)
	}
	log.Printf("Nodes %#v\n", nodes)
	http.Handle("/", _websocket.PingPongConnection(router, hub))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", 10000+*nodeID), nil)
}

func connectToNodes(nodes []string, wallet wallet.Wallet, router _websocket.Router, hub _websocket.Hub) error {
	for _, node := range nodes {
		i, err := strconv.Atoi(node)
		if err != nil {
			return err
		}
		u := url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("localhost:%d", 10000+i),
			Path:   "/",
		}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			return err
		}
		_, err = operations.Register(conn, wallet)(node)
		if err != nil {
			return err
		}
		go _websocket.MaintainConnection(conn, router, hub, node)
	}
	return nil
}
