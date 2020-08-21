package httputil

import (
	"net/http"
)

type redirect struct {
	Redirect bool   `json:"redirect"`
	Code     int    `json:"code"`
	To       string `json:"to"`
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	if ExpectsJson(r) {
		SendJSON(w, redirect{
			Redirect: true,
			Code:     code,
			To:       url,
		}, http.StatusMultipleChoices)
	} else {
		http.Redirect(w, r, url, code)
	}
}
