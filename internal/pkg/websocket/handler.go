package websocket

import (
	"errors"
	"fmt"
)

type Authorizer func(Request) error

type ErrUnauthorized string

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("Node with address %s is unauthorized", string(e))
}

type Handler func(Request) (Response, error)

func (h Handler) Authorized(a Authorizer) Handler {
	return func(request Request) (Response, error) {
		unauthotizedErr := ErrUnauthorized("")
		switch err := a(request); {
		case errors.As(err, &unauthotizedErr):
			return Response{Error: NewUnauthorizedError(err)}, nil
		case err != nil:
			return Response{}, err
		default:
			return h(request)
		}
	}
}

type Router map[CommandType]Handler

func (r Router) Route(c Command) (Response, error) {
	handler, ok := r[c.Type]
	if !ok {
		return Response{Error: NewInvalidCommandError(c.Type)}, nil
	}
	request := Request{
		Body:      c.Body,
		Sender:    c.Sender,
		Signature: c.Signature,
	}
	result, err := handler(request)
	if err != nil {
		return Response{Error: NewUnknownError()}, nil
	}
	return result, nil
}
