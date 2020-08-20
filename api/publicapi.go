package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
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

	api.router.Get("/profile-stat/{profile}", handle.PublicProfileStatsByID(api.db, api.cache))
	return api
}

func (a *PublicAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
