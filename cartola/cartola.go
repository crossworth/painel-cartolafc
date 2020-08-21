package cartola

import (
	"compress/flate"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/auth"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/httputil"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/updater"
	"github.com/crossworth/cartola-web-admin/vk"
	"github.com/crossworth/cartola-web-admin/web"
)

type Cartola struct {
	appName         string
	vkClient        *vk.VKClient
	db              *database.PostgreSQL
	session         sessions.Store
	topicUpdater    *updater.TopicUpdater
	cache           *cache.Cache
	router          chi.Router
	superAdmins     []int
	userTypeHandler *auth.UserTypeHandler
}

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

func NewCartola(
	appName string,
	vkClient *vk.VKClient,
	db *database.PostgreSQL,
	session sessions.Store,
	cache *cache.Cache,
	topicUpdater *updater.TopicUpdater,
	superAdmins []int) *Cartola {
	c := &Cartola{
		appName:      appName,
		vkClient:     vkClient,
		db:           db,
		session:      session,
		cache:        cache,
		topicUpdater: topicUpdater,
		superAdmins:  superAdmins,
	}

	c.userTypeHandler = auth.NewUserTypeHandler(c.db, c.superAdmins)

	c.router = chi.NewRouter()
	c.router.Use(middleware.Recoverer)
	c.router.Use(middleware.Compress(flate.DefaultCompression))
	c.router.Use(cors.New(corsOpts).Handler)
	c.router.Use(middleware.Timeout(10 * time.Minute))
	c.router.Use(middleware.RedirectSlashes)
	c.router.Use(httputil.RemoveDoubleSlashes)

	logger.Log.Info().Msg("montando endpoints")
	c.router.Get("/fazer-login", auth.LoginPage(appName))
	c.router.Get("/login", auth.Login())
	c.router.Get("/login/callback", auth.LoginCallback(c.vkClient, c.session))
	c.router.Get("/logout", auth.Logout(c.session))

	c.router.Mount("/public/api", api.NewPublicAPI(c.db, c.cache))

	c.router.Group(func(r chi.Router) {
		r.Use(auth.OnlyAuthenticatedUsers(c.session, c.userTypeHandler))

		r.Get("/userinfo.js", auth.UserInfo())
		r.Mount("/api", api.NewServer(c.vkClient, c.db, c.cache))
		r.Mount("/", web.New())
	})

	return c
}

func (c *Cartola) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}
