package main

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/auth"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
	"github.com/crossworth/cartola-web-admin/vk/openid"
	"github.com/crossworth/cartola-web-admin/web"
)

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("erro carregando arquivo .env, %v", err)
	}

	logger.Setup(logger.LogInfo, "cartola_web_admin.log")

	logger.Log.Info().Msg("iniciando servidor")

	router := setupRouter()
	vkClient := setupVKClient()
	db := setupDatabase()
	session := setupSessionStorage()
	globalCache := cache.NewCache()

	router.Mount("/public/api", api.NewPublicAPI(db, globalCache))

	setupRoutes(vkClient, session, router, func(router chi.Router) {
		logger.Log.Info().Msg("montando endpoints")
		router.Mount("/api", api.NewServer(vkClient, db, globalCache))
		router.Mount("/", web.New())
	})

	logger.Log.Info().Msgf("iniciando o servidor na porta %s", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router)
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao iniciar o servidor http, %v", err)
	}
}

func setupRoutes(vkAPI *vk.VKClient, sessionStorage sessions.Store, router *chi.Mux, authRoutes func(r chi.Router)) {
	appName := util.GetStringFromEnvOrFatalError("APP_NAME")

	vkAppID := util.GetStringFromEnvOrFatalError("VK_APP_ID")
	vkSecureKey := util.GetStringFromEnvOrFatalError("VK_SECURE_KEY")
	vkCallBackURL := util.GetStringFromEnvOrFatalError("VK_CALLBACK_URL")

	goth.UseProviders(
		openid.New(vkAppID, vkSecureKey, vkCallBackURL, "groups"),
	)

	gothic.Store = sessionStorage

	router.Get("/fazer-login", auth.LoginPage(appName))

	router.Get("/login", auth.Login())
	router.Get("/login/callback", auth.LoginCallback(vkAPI, sessionStorage))
	router.Get("/logout", auth.Logout(sessionStorage))

	router.Group(func(r chi.Router) {
		r.Use(auth.OnlyAuthenticatedUsers(sessionStorage))
		authRoutes(r)
	})
}

func setupSessionStorage() sessions.Store {
	appURL := util.GetStringFromEnvOrFatalError("APP_VK_URL")
	sessionSecret := util.GetStringFromEnvOrFatalError("SESSION_SECRET")

	maxAge := 86400 * 1 // 1 day
	isHttps := strings.HasPrefix(appURL, "https://")

	cookieStore := sessions.NewCookieStore([]byte(sessionSecret))
	cookieStore.MaxAge(maxAge)
	cookieStore.Options.Path = "/"
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.Secure = isHttps
	return cookieStore
}

func setupDatabase() *database.PostgreSQL {
	dbUser := util.GetStringFromEnvOrFatalError("DB_USER")
	dbPassword := util.GetStringFromEnvOrFatalError("DB_PASSWORD")
	dbHost := util.GetStringFromEnvOrFatalError("DB_HOST")
	dbDatabase := util.GetStringFromEnvOrFatalError("DB_DATABASE")

	logger.Log.Info().Msg("conectando no banco de dados")
	db, err := database.NewPostgreSQL(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbDatabase,
	))
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao conectar ao banco de dados, %v", err)
	}
	return db
}

func setupVKClient() *vk.VKClient {
	vkEmail := util.GetStringFromEnvOrFatalError("VK_EMAIL")
	vkPassword := util.GetStringFromEnvOrFatalError("VK_PASSWORD")

	logger.Log.Info().Msg("criando cliente VK")
	vkClient, err := vk.NewVKClient(vkEmail, vkPassword)
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao criar o cliente VK, %v", err)
	}
	return vkClient
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(flate.DefaultCompression))
	router.Use(cors.New(corsOpts).Handler)
	router.Use(middleware.Timeout(10 * time.Minute))
	router.Use(middleware.RedirectSlashes)
	return router
}
