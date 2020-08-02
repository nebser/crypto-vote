package websocket

type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}
