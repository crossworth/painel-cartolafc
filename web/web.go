package web

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr/v2"
)

type Server struct {
	Handler *chi.Mux
}

func New() *Server {
	server := Server{}

	var localFs http.FileSystem
	var assetsFS http.FileSystem

	if _, err := os.Stat("web/frontend/build"); !os.IsNotExist(err) {
		localFs = http.Dir("web/frontend/build")
	} else {
		localFs = packr.New("frontend", "./frontend/build")
	}

	if _, err := os.Stat("web/assets"); !os.IsNotExist(err) {
		assetsFS = http.Dir("web/assets")
	} else {
		assetsFS = packr.New("assets", "./assets")
	}

	fs := http.FileServer(localFs)
	assets := http.FileServer(assetsFS)

	server.Handler = chi.NewRouter()

	server.Handler.HandleFunc("/hahaha.gif", func(writer http.ResponseWriter, request *http.Request) {
		assets.ServeHTTP(writer, request)
	})

	server.Handler.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		assets.ServeHTTP(writer, request)
	})

	server.Handler.HandleFunc("/magicword.mp3", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", " audio/mpeg")
		assets.ServeHTTP(writer, request)
	})

	server.Handler.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
		path := chi.URLParam(request, "*")

		if _, err := localFs.Open(path); os.IsNotExist(err) {
			request.URL.Path = "/"
		}

		fs.ServeHTTP(writer, request)
	})
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}
