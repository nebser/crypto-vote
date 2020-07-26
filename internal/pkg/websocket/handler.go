package websocket

import (
	"log"
	"net/http"
)

type Handler func(resp http.ResponseWriter, request *http.Request) error

func (h Handler) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	if err := h(resp, request); err != nil {
		log.Printf("Error occurred %s\n", err)
	}
}
