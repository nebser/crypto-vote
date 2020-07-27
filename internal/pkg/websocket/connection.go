package websocket

import (
	"log"
	"net/http"
)

type Connection func(resp http.ResponseWriter, request *http.Request) error

func (c Connection) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	if err := c(resp, request); err != nil {
		log.Printf("Error occurred %s\n", err)
	}
}
