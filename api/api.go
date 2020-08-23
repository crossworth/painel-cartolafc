package api

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/httputil"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/vk"
)

type Server struct {
	router chi.Router
	vk     *vk.VKClient
	db     *database.PostgreSQL
	cache  *cache.Cache
}

func NewServer(vk *vk.VKClient, db *database.PostgreSQL, cache *cache.Cache, botQuoteID int) *Server {
	s := &Server{
		router: chi.NewRouter(),
		vk:     vk,
		db:     db,
		cache:  cache,
	}

	logger.SetupLoggerOnRouter(s.router)

	s.router.Use(middleware.RequestID)
	s.router.Use(httputil.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.NoCache)
	s.router.Use(middleware.Compress(flate.DefaultCompression))

	s.router.NotFound(httputil.NotFoundHandler)
	s.router.MethodNotAllowed(httputil.MethodNotAllowedHandler)

	s.router.Get("/my-profile", handle.MyProfile(s.db, s.cache))
	s.router.Get("/my-profile/bot-quotes", handle.MyProfileBotQuotes(s.db, s.cache, botQuoteID))
	s.router.Get("/search", handle.Search(s.db, s.cache))
	s.router.Get("/home-page", handle.GetHomePage(s.db))
	s.router.Get("/topics-ranking", handle.TopicsWithStats(s.db, s.cache))

	// SUPER-ADMIN ONLY
	s.router.Group(func(r chi.Router) {
		r.Use(OnlySuperAdmin())
		r.Get("/administrators-profiles", handle.GetAdministratorProfiles(s.db))
		r.Post("/set-administrators-profiles", handle.SetAdministratorProfiles(s.db))
	})

	// ADMIN ONLY
	s.router.Group(func(r chi.Router) {
		r.Use(OnlyAdmin())

		r.Get("/members-rule", handle.GetMembersRule(s.db))
		r.Post("/set-members-rule", handle.SetMembersRule(s.db))

		r.Post("/set-home-page", handle.SetHomePage(s.db))

		r.Get("/resolve-profile", handle.ProfileLinkToID(s.vk))
		r.Get("/auto-complete/profile/{profile}", handle.AutoCompleteProfileName(s.db))

		r.Get("/profiles", handle.Profiles(s.db, s.cache))
		r.Get("/profiles/{profile}", handle.ProfileByID(s.db))
		r.Get("/profiles/{profile}/history", handle.ProfileHistoryByID(s.db))
		r.Get("/profiles/{profile}/stats", handle.ProfileStatsByID(s.db))
		r.Get("/profiles/{profile}/topics", handle.TopicsByProfileID(s.db))
		r.Get("/profiles/{profile}/comments", handle.CommentsByProfileID(s.db))

		r.Get("/topics", handle.Topics(s.db))
		r.Get("/topics/{topic}", handle.TopicByID(s.db))
		r.Get("/topics/{topic}/comments", handle.CommentFromTopicByID(s.db))
	})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
