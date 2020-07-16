package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/vk"
)

type Server struct {
	router chi.Router
	vk      *vk.VKClient
}

func NewServer(vk *vk.VKClient) *Server {
	s := &Server{}
	s.vk = vk
	s.router = chi.NewRouter()

	s.router.Get("/resolve-profile", handle.ProfileLinkToID(vk))

	// authenticated routes
	s.router.Group(func(r chi.Router) {

	})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
