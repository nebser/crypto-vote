package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
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
		Value:     10,
	}
	return json.Marshal(data)
}

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

func process(wallets wallet.Wallets, parties party.Parties, wg *sync.WaitGroup) error {
	defer wg.Done()
	url := "http://localhost:8000/vote"

	for _, w := range wallets {
		elected := parties[rand.Intn(len(parties))]
		electedPKey := wallet.ExtractPublicKeyHash(elected.Address)
		body := body{
			Sender:    base64.StdEncoding.EncodeToString(w.PublicKeyHash()),
			Recipient: base64.StdEncoding.EncodeToString(electedPKey),
			Verifier:  base64.StdEncoding.EncodeToString(w.PublicKey),
		}
		signature, err := wallet.Sign(body, w.PrivateKey)
		if err != nil {
			return errors.Wrapf(err, "Failed to sign request for %#v", body)
		}
		body.Signature = base64.StdEncoding.EncodeToString(signature)
		raw, err := json.Marshal(body)
		if err != nil {
			return errors.Wrapf(err, "Failed to marshal body %#v", body)
		}
		reader := bytes.NewReader(raw)
		_, err = http.Post(url, "application/json", reader)
		if err != nil {
			return errors.Wrap(err, "Failed to vote")
		}
		log.Printf("Voting for %s\n", elected.Name)
		time.Sleep(2 * time.Second)
	}
	return nil
}

func listParties() (party.Parties, error) {
	response, err := http.Get("http://localhost:8000/parties")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve parties")
	}
	defer response.Body.Close()
	raw, err := ioutil.ReadAll(response.Body)
	var parties party.Parties
	if err := json.Unmarshal(raw, &parties); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal response %s", raw)
	}
	return parties, nil
}

func main() {
	clientKeysDir := flag.String("clients", "clients", "Client key pair files directory")
	flag.Parse()
	files, err := getKeyFiles(*clientKeysDir)
	if err != nil {
		log.Fatalf("Failed to import keys %s", err)
	}
	wallets, err := wallet.ImportMultiple(files)
	if err != nil {
		log.Fatalf("Failed to import wallets %s", err)
	}
	parties, err := listParties()
	if err != nil {
		log.Fatalf("Failed to list parties %s", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := process(wallets, parties, &wg); err != nil {
			log.Printf("Error occurred %s", err)
		}
	}()
	wg.Wait()
}
