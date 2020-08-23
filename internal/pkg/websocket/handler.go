package websocket

import (
	"errors"
	"fmt"
)

type Authorizer func(Ping) error

type ErrUnauthorized string

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("Node with address %s is unauthorized", string(e))
}

type Handler func(Ping, string) (*Pong, error)

func (h Handler) Authorized(a Authorizer) Handler {
	return func(ping Ping, id string) (*Pong, error) {
		unauthotizedErr := ErrUnauthorized("")
		switch err := a(ping); {
		case errors.As(err, &unauthotizedErr):
			return &Pong{
				Message: ErrorMessage,
				Body:    NewUnauthorizedError(err),
			}, nil
		case err != nil:
			return nil, err
		default:
			return h(ping, id)
		}
	}
}

type Router map[Message]Handler

func (r Router) Route(p Ping, id string) *Pong {
	handler, ok := r[p.Message]
	if !ok {
		return &Pong{
			Message: ErrorMessage,
			Body:    NewUnknownMessageError(p.Message),
		}
	}
	result, err := handler(p, id)
	if err != nil {
		return &Pong{
			Message: ErrorMessage,
			Body:    NewUnknownError(),
		}
	}
	return result
}
