package httputil

import (
	"net/http"
	"time"
)

type PaginationMeta struct {
	Prev     string     `json:"prev,omitempty"`
	Current  string     `json:"current"`
	Next     string     `json:"next,omitempty"`
	Total    int        `json:"total"`
	CachedAt *time.Time `json:"cached_at,omitempty"`
}

type Pagination struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func SendPagination(writer http.ResponseWriter, v interface{}, status int, meta PaginationMeta) {
	SendJSON(writer, Pagination{
		Data: v,
		Meta: meta,
	}, status)
}
