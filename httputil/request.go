package httputil

import (
	"net/http"
	"strings"
)

func ExpectsJson(r *http.Request) bool {
	return Ajax(r) || WantsJSON(r)
}

func WantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "/json") || strings.Contains(accept, "+json")
}

func Ajax(r *http.Request) bool {
	return IsXmlHttpRequest(r)
}

func IsXmlHttpRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}
