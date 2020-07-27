package handle

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

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
	message := err.Error()
	status := 500

	if errors.Is(err, sql.ErrNoRows) {
		message = "nenhum resultado encontrado"
		status = 404
	}

	errorCode(writer, NewError(message), status)
}

type PaginationMeta struct {
	Prev     string    `json:"prev,omitempty"`
	Current  string    `json:"current"`
	Next     string    `json:"next,omitempty"`
	Total    int       `json:"total"`
	CachedAt time.Time `json:"cached_at,omitempty"`
}

type Pagination struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func pagination(writer http.ResponseWriter, v interface{}, status int, meta PaginationMeta) {
	json(writer, Pagination{
		Data: v,
		Meta: meta,
	}, status)
}
