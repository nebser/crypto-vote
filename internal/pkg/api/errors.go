package api

import (
	"net/http"
)

type Error struct {
	Error ErrorInformation `json:"error"`
}

type ErrorInformation struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func InternalServerErrorResponse() Response {
	return Response{
		Status: http.StatusInternalServerError,
		Body: Error{
			Error: ErrorInformation{
				Message: "Unexpected error occurred",
				Type:    "internal-server-error",
			},
		},
	}
}

func InvalidDataErrorResponse(message string) Response {
	return Response{
		Status: http.StatusBadRequest,
		Body: Error{
			Error: ErrorInformation{
				Message: message,
				Type:    "invalid-data-error",
			},
		},
	}
}

func UnauthorizedErrorResponse(message string) Response {
	return Response{
		Status: http.StatusUnauthorized,
		Body: Error{
			Error: ErrorInformation{
				Message: message,
				Type:    "unauthorized-error",
			},
		},
	}
}

func UserAlreadyVoted() Response {
	return Response{
		Status: http.StatusConflict,
		Body: Error{
			Error: ErrorInformation{
				Message: "User already voted",
				Type:    "user-already-voted",
			},
		},
	}
}
