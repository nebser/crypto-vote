package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
)

func main() {
	wallets := wallet.Wallets{}
	num := 10
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Argument must be an integer")
		}
		num = parsed
	}
	for i := 0; i < num; i++ {
		w, err := wallet.New()
		if err != nil {
			log.Fatal(err)
		}
		wallets = append(wallets, *w)
	}
	for i, w := range wallets {
		if err := w.Export(fmt.Sprintf("wallets/w%d", i)); err != nil {
			log.Fatal(err)
		}
	}

	w, err := wallet.Import("wallets/w0_pub.pem", "wallets/w0.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(w.Address)
}
