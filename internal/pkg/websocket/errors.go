package websocket

import (
	"fmt"
)

type Error struct {
	Message string `json:"message"`
}

func NewInvalidCommandError(commandType CommandType) *Error {
	return &Error{Message: fmt.Sprintf("Invalid command type %s", commandType)}
}

func NewUnknownError() *Error {
	return &Error{Message: "Unknown error occurred"}
}

func NewInvalidDataError(operation string) *Error {
	return &Error{Message: fmt.Sprintf("Invalid values passed for %s operation", operation)}
}
