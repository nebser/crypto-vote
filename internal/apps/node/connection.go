package node

import (
	"net/http"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func Connection(router _websocket.Router) _websocket.Connection {
	return func(resp http.ResponseWriter, request *http.Request) error {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(resp, request, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to open websocket")
		}
		defer conn.Close()

		for {
			var command _websocket.Command
			if err := conn.ReadJSON(&command); err != nil {
				return errors.Wrap(err, "Failed to parse json into command structure")
			}
			if command.Type == _websocket.CloseConnectionCommand {
				return nil
			}
			response, err := router.Route(command)
			if err != nil {
				response = _websocket.Response{
					Error: _websocket.NewUnknownError(),
				}
			}
			conn.WriteJSON(response)
		}
	}
}
