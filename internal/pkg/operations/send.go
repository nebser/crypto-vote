package operations

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func send(conn *websocket.Conn, op operation, result interface{}) error {
	if err := conn.WriteJSON(op); err != nil {
		return errors.Wrapf(err, "Failed to marshal operation into json %#v", op)
	}
	if err := conn.ReadJSON(result); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal response for operation %s into result", op.Type)
	}
	return nil
}
