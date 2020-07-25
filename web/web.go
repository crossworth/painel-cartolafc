package web

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type Server struct {
	Handler *chi.Mux
}

func New() *Server {
	server := Server{}

	dir := http.Dir("web/frontend/build")
	fs := http.FileServer(dir)

	server.Handler = chi.NewRouter()

	server.Handler.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		path := chi.URLParam(request, "*")

		if _, err := dir.Open(path); os.IsNotExist(err) {
			request.URL.Path = "/"
		}

		fs.ServeHTTP(writer, request)
	})
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}
