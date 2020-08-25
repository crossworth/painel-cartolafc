package cartola

import (
	"compress/flate"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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
	appName                     string
	groupID                     int
	vkClient                    *vk.VKClient
	db                          *database.PostgreSQL
	session                     sessions.Store
	topicUpdater                *updater.TopicUpdater
	cache                       *cache.Cache
	router                      chi.Router
	superAdmins                 []int
	userTypeHandler             *auth.UserHandler
	vkWebhookConfirmationString string
	vkWebhookSecret             string
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
	groupID int,
	vkClient *vk.VKClient,
	db *database.PostgreSQL,
	session sessions.Store,
	cache *cache.Cache,
	topicUpdater *updater.TopicUpdater,
	superAdmins []int,
	botQuoteID int,
	vkWebhookConfirmationString string,
	vkWebhookSecret string) *Cartola {
	c := &Cartola{
		appName:                     appName,
		groupID:                     groupID,
		vkClient:                    vkClient,
		db:                          db,
		session:                     session,
		cache:                       cache,
		topicUpdater:                topicUpdater,
		superAdmins:                 superAdmins,
		vkWebhookConfirmationString: vkWebhookConfirmationString,
		vkWebhookSecret:             vkWebhookSecret,
	}

	c.userTypeHandler = auth.NewUserHandler(c.db, c.superAdmins)

	c.router = chi.NewRouter()
	c.router.Use(middleware.Recoverer)
	c.router.Use(httputil.OnlyAllowedHost)
	c.router.Use(middleware.Compress(flate.DefaultCompression))
	c.router.Use(cors.New(corsOpts).Handler)
	c.router.Use(middleware.Timeout(10 * time.Minute))
	c.router.Use(middleware.RedirectSlashes)
	c.router.Use(httputil.RemoveDoubleSlashes)

	logger.Log.Info().Msg("montando endpoints")
	c.router.HandleFunc("/vk-webhook", c.handleVKWebHook)

	c.router.Get("/fazer-login", auth.LoginPage(appName))
	c.router.Get("/login", auth.Login())
	c.router.Get("/login/callback", auth.LoginCallback(c.vkClient, c.session))
	c.router.Get("/logout", auth.Logout(c.session))

	c.router.Mount("/public/api", api.NewPublicAPI(c.db, c.cache))

	c.router.Group(func(r chi.Router) {
		r.Use(auth.OnlyAuthenticatedUsers(c.session, c.userTypeHandler))

		r.Group(func(noCache chi.Router) {
			noCache.Use(middleware.NoCache)
			noCache.Use(middleware.Compress(flate.DefaultCompression))
			noCache.Get("/userinfo.js", auth.UserInfo())
		})

		r.Mount("/api", api.NewServer(c.vkClient, c.db, c.cache, botQuoteID))
		r.Mount("/", web.New())
	})

	// _, _ = scheduler.Every(1).Day().NotImmediately().Run(c.enqueueAllTopicsIDs)
	// c.enqueueAllTopicsIDs() // fixme: remove this

	return c
}

type VKObject struct {
	TopicID int `json:"topic_id"`
}

type VKEvent struct {
	Type   string   `json:"type"`
	Object VKObject `json:"object"`
	Secret string   `json:"secret"`
}

func (c *Cartola) handleVKWebHook(w http.ResponseWriter, r *http.Request) {
	var event VKEvent

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return
	}

	if event.Type == "confirmation" {
		_, _ = w.Write([]byte(c.vkWebhookConfirmationString))
		return
	}

	if c.vkWebhookSecret != "" && event.Secret != c.vkWebhookSecret {
		return
	}

	if event.Type == "board_post_new" || event.Type == "board_post_edit" && event.Object.TopicID != 0 {
		_ = c.topicUpdater.EnqueueTopicID(event.Object.TopicID)
	}

	_, _ = w.Write([]byte("ok"))
}

func (c *Cartola) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

func (c *Cartola) enqueueAllTopicsIDs() {
	logger.Log.Info().Msg("adicionando tópicos a fila de atualização")
	go func() {
		params := url.Values{}
		params.Set("order", "-2")

		skip := 0
		total := 0

		for {
			params.Set("offset", strconv.Itoa(skip))

			topics, err := c.vkClient.GetClient().BoardGetTopics(c.groupID, 100, params)
			if err != nil {
				logger.Log.Warn().Err(err).Msg("erro ao conseguir todos os ids dos tópicos")
				continue
			}

			for _, topic := range topics.Topics {
				_ = c.topicUpdater.EnqueueTopicIDWithPriority(topic.ID, 15)
			}

			total += len(topics.Topics)

			if total >= topics.Count {
				break
			}

			skip += 100
		}

		logger.Log.Info().Msg("adicionado todos os tópicos a fila")
	}()
}
