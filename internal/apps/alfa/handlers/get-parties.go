package handlers

import (
	"net/http"

	"github.com/nebser/crypto-vote/internal/pkg/api"
	"github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/pkg/errors"
)

func GetParties(getParties party.GetPartiesFn) api.Handler {
	return func(request api.Request) (api.Response, error) {
		parties, err := getParties()
		if err != nil {
			return api.Response{}, errors.Wrapf(err, "Failed to retrieve parties %s", err)
		}
		return api.Response{
			Status: http.StatusOK,
			Body:   parties,
		}, nil
	}
}
