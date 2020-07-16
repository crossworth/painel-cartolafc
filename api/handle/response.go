package handle

import (
	"database/sql"
	"errors"
	"net/http"

	json2 "github.com/helloeave/json"
)

func json(writer http.ResponseWriter, v interface{}, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	out, _ := json2.MarshalSafeCollections(v)
	_, _ = writer.Write(out)
}

func errorCode(writer http.ResponseWriter, err error, status int) {
	json(writer, NewError(err.Error()), status)
}

func databaseError(writer http.ResponseWriter, err error) {
	status := 400

	if errors.Is(err, sql.ErrNoRows) {
		status = 404
	}

	errorCode(writer, NewError(err.Error()), status)
}
