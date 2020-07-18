package handle

import (
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, request *http.Request) {
	json(w, NewError("método não permitido"), 405)
}
