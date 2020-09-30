package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Headers http.Header
	Body    []byte
}

type Response struct {
	Status int
	Body   interface{}
}

type Handler func(Request) (Response, error)

func NewHandleFunc(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			res := InternalServerErrorResponse()
			w.WriteHeader(res.Status)
			json.NewEncoder(w).Encode(res.Body)
			return
		}
		request := Request{
			Headers: r.Header,
			Body:    body,
		}
		result, err := h(request)
		if err != nil {
			log.Printf("Unexpected error occurred %s", err)
			res := InternalServerErrorResponse()
			w.WriteHeader(res.Status)
			json.NewEncoder(w).Encode(res.Body)
			return
		}
		w.WriteHeader(result.Status)
		json.NewEncoder(w).Encode(result.Body)
	}
}
