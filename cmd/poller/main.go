package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/pkg/errors"
)

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

func process(wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		parties, err := listParties()
		if err != nil {
			return errors.Wrap(err, "Failed to list parties")
		}
		fmt.Println("START PARTY LIST")
		for _, p := range parties {
			fmt.Printf("%s:\t%d\n", p.Name, p.Balance/10)
		}
		fmt.Println("END PARTY LIST")
		time.Sleep(10 * time.Second)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := process(&wg); err != nil {
			fmt.Printf("Unexpected error occurred %s\n", err)
		}
	}()
	wg.Wait()
}
