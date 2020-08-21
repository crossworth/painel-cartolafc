package httputil

import (
	"database/sql"
	"errors"
	"net/http"
)

func SendDatabaseError(writer http.ResponseWriter, err error) {
	message := err.Error()
	status := 500

	if errors.Is(err, sql.ErrNoRows) {
		message = "nenhum resultado encontrado"
		status = 404
	}

	SendErrorCode(writer, NewError(message), status)
}
