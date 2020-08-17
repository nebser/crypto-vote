package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Connection func(resp http.ResponseWriter, request *http.Request) error

func (c Connection) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	if err := c(resp, request); err != nil {
		log.Printf("Error occurred %s\n", err)
	}
}

func reader(conn *websocket.Conn, router Router, responseChan chan Pong, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(responseChan)
	for {
		var ping Ping
		if err := conn.ReadJSON(&ping); err != nil {
			log.Println("Failed to parse message")
			responseChan <- Pong{
				Message: ErrorMessage,
			}
			continue
		}
		if ping.Message == CloseConnectionMessage {
			return
		}
		pong := router.Route(ping)
		if pong != nil {
			responseChan <- *pong
		}
	}
}

func writer(conn *websocket.Conn, responseChan chan Pong, wg *sync.WaitGroup) {
	defer wg.Done()
	for pong := range responseChan {
		conn.WriteJSON(pong)
	}
}

func PingPongConnection(router Router) Connection {
	return func(resp http.ResponseWriter, request *http.Request) error {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(resp, request, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to open websocket")
		}
		defer conn.Close()

		responseChan := make(chan Pong, 5)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go reader(conn, router, responseChan, &wg)
		go writer(conn, responseChan, &wg)

		wg.Wait()

		return nil
	}
}
