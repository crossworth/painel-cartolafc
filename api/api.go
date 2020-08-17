package api

import (
	"compress/flate"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/hlog"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/vk"
)

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

	s.setupLogger()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.NoCache)
	s.router.Use(middleware.Compress(flate.DefaultCompression))

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
		r.Get("/topics-ranking", handle.TopicsWithStats(s.db, s.cache))
		r.Get("/topics/{topic}", handle.TopicByID(s.db))
		r.Get("/topics/{topic}/comments", handle.CommentFromTopicByID(s.db))

		r.Get("/search", handle.Search(s.db, s.cache))
	})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupLogger() {
	s.router.Use(hlog.NewHandler(logger.Log))

	s.router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		l := hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration)

		if vkID, ok := model.VKIDFromRequest(r); ok {
			l.Str("vk_id", strconv.Itoa(vkID))
		}

		l.Msg("")
	}))

	s.router.Use(hlog.RemoteAddrHandler("ip"))
	s.router.Use(hlog.UserAgentHandler("user_agent"))
	s.router.Use(hlog.RefererHandler("referer"))
	s.router.Use(hlog.RequestIDHandler("req_id", "request-id"))
}
