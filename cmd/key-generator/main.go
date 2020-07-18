package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
)

func main() {
	alfaKeyDir := flag.String("alfa", "alfa", "Directory where to create key pairs for alfa node")
	clientKeysDir := flag.String("clients", "clients", "Directory where to create client key pairs")
	num := flag.Int("clientsNumber", 50, "Number of client key pairs to generate")
	flag.Parse()

	clientWallets := wallet.Wallets{}
	for i := 0; i < *num; i++ {
		w, err := wallet.New()
		if err != nil {
			log.Fatal(err)
		}
		clientWallets = append(clientWallets, *w)
	}
	for i, w := range clientWallets {
		if err := w.Export(fmt.Sprintf("%s/c%d", *clientKeysDir, i)); err != nil {
			log.Fatal(err)
		}
	}

	alfaWallet, err := wallet.New()
	if err != nil {
		log.Fatalf("Failed to create wallet for alfa node. Error %s", err)
	}
	if err := alfaWallet.Export(fmt.Sprintf("%s/key", *alfaKeyDir)); err != nil {
		log.Fatal(err)
	}
}
