package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
)

func exportMultiple(directory, base string, start, num int) error {
	wallets := wallet.Wallets{}
	for i := 0; i < num; i++ {
		w, err := wallet.New()
		if err != nil {
			return err
		}
		wallets = append(wallets, *w)
	}
	for i, w := range wallets {
		if err := w.Export(fmt.Sprintf("%s/%s%d", directory, base, start+i)); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	alfaKeyDir := flag.String("alfa", "alfa", "Directory where to create key pairs for alfa node")
	clientKeysDir := flag.String("clients", "clients", "Directory where to create client key pairs")
	nodesKeysDir := flag.String("nodes", "nodes", "Directory where to create node key pairs")
	numOfClients := flag.Int("clientsNumber", 50, "Number of client key pairs to generate")
	numOfNodes := flag.Int("nodesNumber", 5, "Number of node key pairs to generate")
	flag.Parse()

	if err := exportMultiple(*clientKeysDir, "c", 0, *numOfClients); err != nil {
		log.Fatalf("Failed to generate keys for clients %s", err)
	}
	if err := exportMultiple(*nodesKeysDir, "n", 1, *numOfNodes); err != nil {
		log.Fatalf("Failed to generate keys for nodes %s", err)
	}

	alfaWallet, err := wallet.New()
	if err != nil {
		log.Fatalf("Failed to create wallet for alfa node. Error %s", err)
	}
	if err := alfaWallet.Export(fmt.Sprintf("%s/key", *alfaKeyDir)); err != nil {
		log.Fatal(err)
	}
}
