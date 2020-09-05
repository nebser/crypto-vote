package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	data := struct {
		Error Error `json:"error"`
	}{
		Error: e,
	}
	return json.Marshal(data)
}

func InternalServerErrorResponse() Response {
	return Response{
		Status: http.StatusInternalServerError,
		Body: Error{
			Message: "Unexpected error occurred",
			Type:    "internal-server-error",
		},
	}
}

func InvalidDataErrorResponse(message string) Response {
	return Response{
		Status: http.StatusBadRequest,
		Body: Error{
			Message: message,
			Type:    "invalid-data-error",
		},
	}
}

func UnauthorizedErrorResponse(message string) Response {
	return Response{
		Status: http.StatusUnauthorized,
		Body: Error{
			Message: message,
			Type:    "unauthorized-error",
		},
	}
}
