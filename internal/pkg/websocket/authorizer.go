package websocket

import (
	"fmt"

	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type Authorizer func(Ping) error

type ErrUnauthorized string

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("Node with address %s is unauthorized", string(e))
}

func PublicKeyAuthorizer(pkey string, verify wallet.VerifierFn) Authorizer {
	return func(ping Ping) error {
		switch ok, err := verify(ping, ping.Signature, pkey); {
		case err != nil:
			return errors.Wrap(err, "Failed to verify signature")
		case !ok:
			return ErrUnauthorized("Invalid signature")
		default:
			return nil
		}
	}
}
