package node

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

func reader(conn *websocket.Conn, router _websocket.PingPongRouter, responseChan chan _websocket.Pong, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(responseChan)
	for {
		var ping _websocket.Ping
		if err := conn.ReadJSON(&ping); err != nil {
			log.Println("Failed to parse message")
			responseChan <- _websocket.Pong{
				Message: _websocket.ErrorMessage,
			}
			continue
		}
		if ping.Message == _websocket.CloseConnectionMessage {
			return
		}
		pong := router.Route(ping)
		if pong != nil {
			responseChan <- *pong
		}
	}
}

func writer(conn *websocket.Conn, responseChan chan _websocket.Pong, wg *sync.WaitGroup) {
	defer wg.Done()
	for pong := range responseChan {
		conn.WriteJSON(pong)
	}
}

func PingPongConnection(router _websocket.PingPongRouter) _websocket.Connection {
	return func(resp http.ResponseWriter, request *http.Request) error {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(resp, request, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to open websocket")
		}
		defer conn.Close()

		responseChan := make(chan _websocket.Pong, 5)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go reader(conn, router, responseChan, &wg)
		go writer(conn, responseChan, &wg)

		wg.Wait()

		return nil
	}
}

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
				response = _websocket.Response{}
			}
			conn.WriteJSON(response)
		}
	}
}
