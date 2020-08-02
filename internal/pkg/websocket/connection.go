package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Connection func(resp http.ResponseWriter, request *http.Request) error

func (c Connection) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	if err := c(resp, request); err != nil {
		log.Printf("Error occurred %s\n", err)
	}
}

type Reader func(interface{}) error

type Writer func(interface{}) error

type ConnectionHandler func(Router, Reader, Writer) error

func NewSocketConnection(router Router, handler ConnectionHandler) Connection {
	return func(resp http.ResponseWriter, request *http.Request) error {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(resp, request, nil)
		if err != nil {
			return errors.Wrap(err, "Failed to open websocket")
		}
		defer conn.Close()

		return handler(router, conn.ReadJSON, conn.WriteJSON)
	}
}
