package httputil

import (
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, request *http.Request) {
	SendJSON(w, NewError("método não permitido"), 405)
}

func NotFoundHandler(w http.ResponseWriter, request *http.Request) {
	SendJSON(w, NewError("não encontrado"), 404)
}
