package cartola

import (
	"compress/flate"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"

	"github.com/crossworth/painel-cartolafc/api"
	"github.com/crossworth/painel-cartolafc/auth"
	"github.com/crossworth/painel-cartolafc/cache"
	"github.com/crossworth/painel-cartolafc/database"
	"github.com/crossworth/painel-cartolafc/httputil"
	"github.com/crossworth/painel-cartolafc/logger"
	"github.com/crossworth/painel-cartolafc/updater"
	"github.com/crossworth/painel-cartolafc/vk"
	"github.com/crossworth/painel-cartolafc/web"
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
	c.router.Mount("/debug", middleware.Profiler())

	cartola := chi.NewRouter()
	cartola.Use(middleware.Recoverer)
	cartola.Use(httputil.OnlyAllowedHost)
	cartola.Use(middleware.Compress(flate.DefaultCompression))
	cartola.Use(cors.New(corsOpts).Handler)
	cartola.Use(middleware.Timeout(10 * time.Minute))
	cartola.Use(middleware.RedirectSlashes)
	cartola.Use(httputil.RemoveDoubleSlashes)

	logger.Log.Info().Msg("montando endpoints")
	cartola.HandleFunc("/vk-webhook", c.handleVKWebHook)

	cartola.Get("/fazer-login", auth.LoginPage(appName))
	cartola.Get("/login", auth.Login())
	cartola.Get("/login/callback", auth.LoginCallback(c.vkClient, c.session))
	cartola.Get("/logout", auth.Logout(c.session))

	cartola.Mount("/public/api", api.NewPublicAPI(c.db, c.cache))

	cartola.Group(func(r chi.Router) {
		r.Use(auth.OnlyAuthenticatedUsers(c.session, c.userTypeHandler))

		r.Group(func(noCache chi.Router) {
			noCache.Use(middleware.NoCache)
			noCache.Use(middleware.Compress(flate.DefaultCompression))
			noCache.Get("/userinfo.js", auth.UserInfo())
		})

		r.Mount("/api", api.NewServer(c.vkClient, c.db, c.cache, botQuoteID))
		r.Mount("/", web.New())
	})

	c.router.Mount("/", cartola)

	// NOTE(Pedro): Every 10 days this will run and enqueue all the
	// community topics to update the likes (in a global way)
	// and possible get missed topics
	_, _ = scheduler.Every(10).Day().Run(c.enqueueAllTopicsIDs)

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

			var ids []int
			for _, topic := range topics.Topics {
				ids = append(ids, topic.ID)
			}

			_ = c.topicUpdater.EnqueueTopicSliceIDWithPriority(ids, 15)
			total += len(topics.Topics)

			logger.Log.Info().Msgf("adicionado %d tópicos a fila", total)
			if total >= topics.Count {
				break
			}

			skip += 100
		}

		logger.Log.Info().Msg("adicionado todos os tópicos a fila")
	}()
}
