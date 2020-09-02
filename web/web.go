package web

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr/v2"
)

type Server struct {
	Handler *chi.Mux
}

func New() *Server {
	server := Server{}

	var localFs http.FileSystem

	if _, err := os.Stat("web/frontend/build"); !os.IsNotExist(err) {
		localFs = http.Dir("web/frontend/build")
	} else {
		localFs = packr.New("frontend", "./frontend/build")
	}

	fs := http.FileServer(localFs)

	server.Handler = chi.NewRouter()
	server.Handler.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		path := chi.URLParam(request, "*")

		file, err := localFs.Open(path)
		if os.IsNotExist(err) || strings.HasSuffix(request.URL.Path, "/") {
			request.URL.Path = "/"
		} else {
			_ = file.Close()
		}

		fs.ServeHTTP(writer, request)
	})
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}
