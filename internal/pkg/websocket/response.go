package websocket

type Response struct {
	Data  interface{} `json:"data"`
	Error *Error      `json:"error"`
}
