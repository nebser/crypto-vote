package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/repository"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"

	"github.com/nebser/crypto-vote/internal/apps/alfa"
	"github.com/nebser/crypto-vote/internal/apps/alfa/handlers"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
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

	getTip := repository.GetTip(db)
	getBlock := repository.GetBlock(db)
	saveNode := repository.SaveNode(db)
	if *newOption {
		if err := alfa.Initialize(
			*masterWallet,
			wallets,
			repository.InitBlockchain(db),
			repository.AddBlock(db),
			saveNode); err != nil {
			log.Fatal(err)
		}
	}
	blockchain.PrintBlockchain(getTip, getBlock)
	router := websocket.Router{
		websocket.GetBlockchainHeightCommand: handlers.GetHeightHandler(getTip, getBlock),
		websocket.GetMissingBlocksCommand:    handlers.GetMissingBlocks(getTip, getBlock),
		websocket.GetBlockCommand:            handlers.GetBlock(getBlock),
		websocket.RegisterCommand:            handlers.Register(saveNode, repository.GetNodes(db)),
	}
	http.Handle("/", alfa.Connection(router))
	http.ListenAndServe(":10000", nil)
}
