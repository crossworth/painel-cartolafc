package httputil

import (
	"net/http"
)

type Error struct {
	IsError bool   `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(text string) error {
	return &Error{
		IsError: true,
		Message: text,
	}
}

func SendErrorCode(writer http.ResponseWriter, err error, status int) {
	SendJSON(writer, NewError(err.Error()), status)
}
