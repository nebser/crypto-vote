package websocket

import (
	"fmt"
)

const (
	UnknownMessageErrorName = "message-unknown"
	UnknownErrorName        = "unknown-error"
	UnauthorizedErrorName   = "unauthorized"
	BlockNotFoundErrorName  = "block-not-found"
	InvalidDataErrorName    = "invalid-data"
)

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewUnknownError() Error {
	return Error{
		Name:    UnknownErrorName,
		Message: fmt.Sprintf("Unknown error occurred"),
	}
}

func NewInvalidDataError(operation string) Error {
	return Error{
		Name:    InvalidDataErrorName,
		Message: fmt.Sprintf("Invalid values passed for %s operation", operation),
	}
}

func NewBlockNotFoundError(blockHash []byte) Error {
	return Error{
		Name:    BlockNotFoundErrorName,
		Message: fmt.Sprintf("Block %x not found", blockHash),
	}
}

func NewUnauthorizedError(err error) Error {
	return Error{
		Name:    UnauthorizedErrorName,
		Message: fmt.Sprintf("Unathorized. Error: %s", err),
	}
}

func NewUnknownMessageError(message Message) Error {
	return Error{
		Name:    UnknownMessageErrorName,
		Message: fmt.Sprintf("Unknown message %s", message),
	}
}
