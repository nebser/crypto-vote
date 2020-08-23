package websocket

import (
	"io"
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

func reader(conn *websocket.Conn, id string, hub Hub, router Router, responseChan chan Pong, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(responseChan)
	defer hub.Unregister(id)
	for {
		var ping Ping
		if err := conn.ReadJSON(&ping); err != nil {
			if err != io.ErrUnexpectedEOF {
				log.Println("Closing reader")
				return
			}
			log.Printf("Failed to parse message %+v, %t\n", err, errors.Is(err, io.ErrUnexpectedEOF))
			responseChan <- Pong{
				Message: ErrorMessage,
			}
			continue
		}
		if ping.Message == CloseConnectionMessage {
			return
		}
		pong := router.Route(ping, id)
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

func PingPongConnection(router Router, hub Hub) Connection {
	return func(resp http.ResponseWriter, request *http.Request) error {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(resp, request, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to open websocket")
		}
		defer conn.Close()

		responseChan := make(chan Pong, 5)
		id := hub.Add(responseChan)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go reader(conn, id, hub, router, responseChan, &wg)
		go writer(conn, responseChan, &wg)

		wg.Wait()

		return nil
	}
}
