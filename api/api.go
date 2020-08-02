package api

import (
	"compress/flate"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/mattn/go-colorable"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/vk"
)

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

type Server struct {
	router chi.Router
	vk     *vk.VKClient
	db     *database.PostgreSQL
	cache  *cache.Cache
}

func NewServer(vk *vk.VKClient, db *database.PostgreSQL) *Server {
	s := &Server{
		router: chi.NewRouter(),
		vk:     vk,
		db:     db,
		cache:  cache.NewCache(),
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.NoCache)
	s.router.Use(middleware.Compress(flate.DefaultCompression))
	s.router.Use(cors.New(corsOpts).Handler)
	s.router.Use(middleware.RequestLogger(getLogger()))

	// TODO(Pedro): add timeout to routes?

	s.router.NotFound(handle.NotFoundHandler)
	s.router.MethodNotAllowed(handle.MethodNotAllowedHandler)

	// public routes
	s.router.Get("/resolve-profile", handle.ProfileLinkToID(s.vk))

	// authenticated routes
	s.router.Group(func(r chi.Router) {
		r.Get("/auto-complete/profile/{profile}", handle.AutoCompleteProfileName(s.db))

		r.Get("/profiles", handle.Profiles(s.db, s.cache))
		r.Get("/profiles/{profile}", handle.ProfileByID(s.db))
		r.Get("/profiles/{profile}/history", handle.ProfileHistoryByID(s.db))
		r.Get("/profiles/{profile}/stats", handle.ProfileStatsByID(s.db))
		r.Get("/profiles/{profile}/topics", handle.TopicsByProfileID(s.db))
		r.Get("/profiles/{profile}/comments", handle.CommentsByProfileID(s.db))

		r.Get("/topics", handle.Topics(s.db))
		// r.Get("/topics/{topic}", handle.TopicByID(s.db))
		// r.Get("/topics/{topic}/comments", handle.TopicByID(s.db))
	})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func getLogger() middleware.LogFormatter {
	var logOutput io.Writer
	logOutput = os.Stdout

	if runtime.GOOS == "windows" {
		logOutput = colorable.NewColorableStdout()
	}

	return &middleware.DefaultLogFormatter{Logger: log.New(logOutput, "", log.LstdFlags), NoColor: false}
}
