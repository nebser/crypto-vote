package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
)

type body struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Verifier  string `json:"verifier"`
	Signature string `json:"signature"`
}

func (b body) Signable() ([]byte, error) {
	data := struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Value     int    `json:"value"`
	}{
		Sender:    b.Sender,
		Recipient: b.Recipient,
		Value:     1,
	}
	return json.Marshal(data)
}

func main() {
	url := "http://localhost:8000/vote"
	id := flag.Int("id", -1, "ID of the client that's voting")
	choice := flag.Int("choice", -1, "ID of the choice to vote for")
	flag.Parse()
	if *id == -1 {
		log.Fatalf("ID flag must be greater or equal to zero")
	}
	if *choice == -1 {
		log.Fatalf("Choice flag must be greater or equal to zero")
	}
	keyfiles := keyfiles.KeyFiles{
		PrivateKeyFile: fmt.Sprintf("clients/c%d.pem", *id),
		PublicKeyFile:  fmt.Sprintf("clients/c%d_pub.pem", *id),
	}
	w, err := wallet.Import(keyfiles)
	if err != nil {
		panic(err)
	}
	partyPub, err := wallet.LoadPublicKey(fmt.Sprintf("nodes/n%d_pub.pem", *choice))
	if err != nil {
		panic(err)
	}
	hashedPartyPub, err := wallet.HashedPublicKey(partyPub)
	if err != nil {
		panic(err)
	}
	body := body{
		Sender:    base64.StdEncoding.EncodeToString(w.PublicKeyHash()),
		Recipient: base64.StdEncoding.EncodeToString(hashedPartyPub),
		Verifier:  base64.StdEncoding.EncodeToString(w.PublicKey),
	}
	signature, err := wallet.Sign(body, w.PrivateKey)
	if err != nil {
		panic(err)
	}
	body.Signature = base64.StdEncoding.EncodeToString(signature)
	raw, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(raw)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("Received response %s", result)
}
