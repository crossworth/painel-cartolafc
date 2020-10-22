package api

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crossworth/painel-cartolafc/api/handle"
	"github.com/crossworth/painel-cartolafc/cache"
	"github.com/crossworth/painel-cartolafc/database"
	"github.com/crossworth/painel-cartolafc/logger"
)

type PublicAPI struct {
	router chi.Router
	db     *database.PostgreSQL
	cache  *cache.Cache
}

func NewPublicAPI(db *database.PostgreSQL, cache *cache.Cache) *PublicAPI {
	api := &PublicAPI{
		router: chi.NewRouter(),
		db:     db,
		cache:  cache,
	}

	logger.SetupLoggerOnRouter(api.router)

	api.router.Use(middleware.NoCache)
	api.router.Use(middleware.Compress(flate.DefaultCompression))
	api.router.Get("/profile-stat/{profile}", handle.PublicProfileStatsByID(api.db, api.cache))
	return api
}

func (a *PublicAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
