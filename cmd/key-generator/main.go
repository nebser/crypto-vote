package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
)

func main() {
	wallets := wallet.Wallets{}
	num := 50
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Argument must be an integer")
			os.Exit(1)
		}
		num = parsed
	}
	for i := 0; i < num; i++ {
		w, err := wallet.New()
		if err != nil {
			fmt.Println(err)
			return
		}
		wallets = append(wallets, *w)
	}
	serlialized, err := wallets.Serialized()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := ioutil.WriteFile("wallets.json", serlialized, 0644); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SUCCESS")
}
