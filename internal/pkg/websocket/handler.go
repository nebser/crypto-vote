package websocket

import "encoding/json"

type Handler func(payload json.RawMessage) (Response, error)

type Router map[CommandType]Handler

func (r Router) Route(c Command) (Response, error) {
	handler, ok := r[c.Type]
	if !ok {
		return Response{Error: NewInvalidCommandError(c.Type)}, nil
	}
	return handler(c.Body)
}
