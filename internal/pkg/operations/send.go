package operations

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"
	"github.com/pkg/errors"
)

type response struct {
	Result json.RawMessage   `json:"result"`
	Error  *_websocket.Error `json:"error"`
}

func send(conn *websocket.Conn, op operation, result interface{}) error {
	if err := conn.WriteJSON(op); err != nil {
		return errors.Wrapf(err, "Failed to marshal operation into json %#v", op)
	}
	if err := conn.ReadJSON(result); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal response for operation %s into result", op.Type)
	}
	return nil
}

func call(conn *websocket.Conn, op operation, result interface{}) error {
	var r response
	if err := send(conn, op, &r); err != nil {
		return errors.Wrapf(err, "Failed to send operation %#v", op)
	}
	if r.Error != nil {
		return errors.Errorf("Failed to perform operation %#v. Error: %s", op, r.Error)
	}
	if err := json.Unmarshal(r.Result, result); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal response %s", r.Result)
	}
	return nil
}
