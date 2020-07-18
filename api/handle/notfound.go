package handle

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, request *http.Request) {
	json(w, NewError("não encontrado"), 404)
}
