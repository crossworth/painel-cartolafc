package httputil

import (
	"net/http"

	json2 "github.com/helloeave/json"
)

func SendJSON(writer http.ResponseWriter, v interface{}, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	out, _ := json2.MarshalSafeCollections(v)
	_, _ = writer.Write(out)
}
