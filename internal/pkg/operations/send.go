package operations

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type response struct {
	Message   _websocket.Message `json:"message"`
	Body      json.RawMessage    `json:"body"`
	Signature string             `json:"signature"`
	Sender    string             `json:"sender"`
}

func send(conn *websocket.Conn, op operation, result interface{}) error {
	if err := conn.WriteJSON(op); err != nil {
		return errors.Wrapf(err, "Failed to marshal operation into json %#v", op)
	}
	if err := conn.ReadJSON(result); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal response for operation %s into result", op.Message)
	}
	return nil
}

func call(conn *websocket.Conn, op operation, result interface{}) error {
	var r response
	if err := send(conn, op, &r); err != nil {
		return errors.Wrapf(err, "Failed to send operation %#v", op)
	}
	if r.Message == _websocket.ErrorMessage {
		return errors.Errorf("Failed to perform operation %#v. Error: %s", op, r.Body)
	}
	if err := json.Unmarshal(r.Body, result); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal response %s", r.Body)
	}
	return nil
}
