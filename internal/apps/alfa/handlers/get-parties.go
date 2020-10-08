package handlers

import (
	"net/http"
	"sort"

	"github.com/nebser/crypto-vote/internal/pkg/api"
	"github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

func GetParties(getParties party.GetPartiesFn, getUTXOsByPublicKey transaction.GetUTXOsByPublicKeyFn) api.Handler {
	return func(request api.Request) (api.Response, error) {
		parties, err := getParties()
		if err != nil {
			return api.Response{}, errors.Wrapf(err, "Failed to retrieve parties %s", err)
		}
		result := make(party.Parties, 0, cap(parties))
		for _, p := range parties {
			utxos, err := getUTXOsByPublicKey(wallet.ExtractPublicKeyHash(p.Address))
			if err != nil {
				return api.Response{}, errors.Wrapf(err, "Failed to enrich party with balance %#v", p)
			}
			enriched := p
			enriched.Balance = utxos.Sum()
			result = append(result, enriched)
		}
		sort.Sort(sort.Reverse(result))
		return api.Response{
			Status: http.StatusOK,
			Body:   result,
		}, nil
	}
}
