package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nebser/crypto-vote/internal/pkg/api"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"

	"github.com/gorilla/mux"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/repository"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"

	"github.com/nebser/crypto-vote/internal/apps/alfa"
	"github.com/nebser/crypto-vote/internal/apps/alfa/handlers"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/robfig/cron/v3"
)

const (
	dbFileName = "db"
)

func getKeyFiles(keyDirectory string) (keyfiles.KeyFilesList, error) {
	files, err := ioutil.ReadDir(keyDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read key file directory %s", keyDirectory)
	}

	fileGroups := map[string]keyfiles.KeyFiles{}
	for _, f := range files {
		if strings.Contains(f.Name(), "address") {
			continue
		}
		name := strings.Replace(f.Name(), "_pub", "", 1)
		group := fileGroups[name]
		if strings.Contains(f.Name(), "pub") {
			group.PublicKeyFile = fmt.Sprintf("%s/%s", keyDirectory, f.Name())
		} else {
			group.PrivateKeyFile = fmt.Sprintf("%s/%s", keyDirectory, f.Name())
		}
		fileGroups[name] = group
	}

	result := keyfiles.KeyFilesList{}
	for _, keyFiles := range fileGroups {
		result = append(result, keyFiles)
	}
	return result, nil
}

func main() {
	newOption := flag.Bool("new", false, "Should initialize new blockchain")
	privateKey := flag.String("private", "alfa/key.pem", "Private key file path")
	publicKey := flag.String("public", "alfa/key_pub.pem", "Public key file path")
	clientKeysDir := flag.String("clients", "clients", "Client key pair files directory")
	nodeKeysDir := flag.String("nodes", "nodes", "Nodes key pair files directory")

	flag.Parse()
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
	masterWallet, err := wallet.Import(keyfiles.KeyFiles{
		PublicKeyFile:  *publicKey,
		PrivateKeyFile: *privateKey,
	})
	if err != nil {
		log.Fatalf("Failed to load master wallet %s", err)
	}
	clientKeyFiles, err := getKeyFiles(*clientKeysDir)
	if err != nil {
		log.Fatalf("Failed to load client key files directory %s", err)
	}
	nodeKeyFiles, err := getKeyFiles(*nodeKeysDir)
	if err != nil {
		log.Fatalf("Failed to load node key files directory %s", err)
	}
	wallets, err := wallet.ImportMultiple(append(clientKeyFiles, nodeKeyFiles...))
	if err != nil {
		log.Fatalf("Failed to wallets %s", err)
	}

	if *newOption {
		if err := alfa.Initialize(
			*masterWallet,
			wallets,
			repository.InitBlockchain(db),
			repository.AddBlock(db)); err != nil {
			log.Fatal(err)
		}
	}
	blockchain.PrintBlockchain(repository.GetTip(db), repository.GetBlock(db))
	hub := websocket.NewHub()
	startForgerChooser(db, *masterWallet, hub)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go runSocketServer(&wg, db, hub, *masterWallet)
	go runAPIServer(&wg, db, hub)
	wg.Wait()
}

func startForgerChooser(db *bolt.DB, masterWallet wallet.Wallet, hub *websocket.Hub) {
	c := cron.New()
	c.Schedule(
		cron.Every(30*time.Second),
		alfa.Runner(
			hub.RegisteredNodes,
			hub.RandomUnicast,
			repository.GetTip(db),
			repository.GetBlock(db),
			wallet.NewSigner(masterWallet),
		),
	)
	c.Start()
}

func runSocketServer(wg *sync.WaitGroup, db *bolt.DB, hub *websocket.Hub, w wallet.Wallet) {
	defer wg.Done()
	getTip := repository.GetTip(db)
	getBlock := repository.GetBlock(db)
	findBlock := blockchain.FindBlock(getTip, getBlock)
	authorizer := blockchain.BlockchainAuthorizer(findBlock)
	router := websocket.Router{
		websocket.GetBlockchainHeightMessage: handlers.GetHeightHandler(getTip, getBlock),
		websocket.GetMissingBlocksMessage:    handlers.GetMissingBlocks(getTip, getBlock),
		websocket.GetBlockMessage:            handlers.GetBlock(getBlock),
		websocket.RegisterMessage:            handlers.Register(hub).Authorized(authorizer),
		websocket.BlockForgedMessage: handlers.BlockForged(
			getTip,
			getBlock,
			blockchain.VerfiyBlock(
				transaction.VerifyTransactions(
					repository.GetTransactionUTXO(db),
					wallet.VerifySignature,
				),
				transaction.IsStakeTransaction(w.PublicKeyHash()),
			),
			repository.AddNewBlock(db),
		),
	}
	mux := http.NewServeMux()
	mux.Handle("/", websocket.PingPongConnection(router, hub, wallet.NewSigner(w)))
	http.ListenAndServe(":10000", mux)
}

func runAPIServer(wg *sync.WaitGroup, db *bolt.DB, hub *websocket.Hub) {
	getTip := repository.GetTip(db)
	getBlock := repository.GetBlock(db)
	findBlock := blockchain.FindBlock(getTip, getBlock)
	httpRouter := mux.NewRouter()
	httpRouter.
		HandleFunc("/vote",
			api.NewHandleFunc(
				handlers.Vote(
					findBlock,
					repository.CastVote(db),
					hub.Broadcast,
				),
			),
		).Methods("POST")
	serverMux := http.NewServeMux()
	serverMux.Handle("/", httpRouter)
	http.ListenAndServe(":8000", serverMux)
}
