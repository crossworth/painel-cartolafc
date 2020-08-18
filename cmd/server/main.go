package main

import (
	"compress/flate"
	"context"
	"fmt"
	"html/template"
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
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
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
	setupRoutes(vkClient, session, router, func(router chi.Router) {
		logger.Log.Info().Msg("montando endpoints")
		router.Mount("/api", api.NewServer(vkClient, db))
		router.Mount("/", web.New())
	})

	logger.Log.Info().Msgf("iniciando o servidor na porta %s\n", os.Getenv("APP_PORT"))

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

	loginPageTemplate, err := template.New("loginPage").Parse(loginPage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("erro ao criar template de login")
	}

	hahahaPage, err := template.New("hahahaPage").Parse(hahahaPage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("erro ao criar template de hahahaPage")
	}

	hahahaHandler := func(message string) func(writer http.ResponseWriter, request *http.Request) {
		return func(writer http.ResponseWriter, request *http.Request) {
			hahahaPage.Execute(writer, struct {
				Message string
			}{
				Message: message,
			})
		}
	}

	router.Get("/hahaha", hahahaHandler(""))

	router.Get("/fazer-login", func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		_, err := gothic.CompleteUserAuth(writer, request)
		if err == nil {
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			return
		}

		loginPageTemplate.Execute(writer, struct {
			Title string
		}{
			Title: appName,
		})
	})

	router.Get("/login", func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		gothic.BeginAuthHandler(writer, request)
	})

	router.Get("/login/callback", func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		user, err := gothic.CompleteUserAuth(writer, request)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao fazer segunda parte do login")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		session, err := sessionStorage.Get(request, model.UserSession)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao conseguir a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		session.Values["user_id"] = user.UserID
		// todo: check user is on group
		isOnCommunity, err := vkAPI.IsUserIDOnGroup(request.Context(), user.UserID,
			util.GetIntFromEnvOrFatalError("APP_VK_GROUP_ID"))
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao salvar a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		if !isOnCommunity {
			_ = gothic.Logout(writer, request)
			hahahaHandler("Você não é um membro da comunidade")(writer, request)
			return
		}

		err = session.Save(request, writer)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao salvar a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
	})

	router.Get("/logout", func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))

		session, _ := sessionStorage.Get(request, model.UserSession)
		delete(session.Values, "user_id")
		_ = session.Save(request, writer)

		_ = gothic.Logout(writer, request)
		http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
	})

	router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.String() == "/hahaha.gif" || request.URL.String() == "/favicon.ico" ||
					request.URL.String() == "/magicword.mp3" {
					next.ServeHTTP(writer, request)
					return
				}

				session, err := sessionStorage.Get(request, model.UserSession)
				if err != nil {
					logger.Log.Error().Err(err).Msg("erro ao conseguir a session de usuário no middleware")
					hahahaHandler("Ocorreu um problema")(writer, request)
					return
				}

				userID, ok := session.Values["user_id"].(string)
				if !ok {
					hahahaHandler("Usuário não logado")(writer, request)
					return
				}

				// todo(pedro): check time?
				if userID != "" {
					next.ServeHTTP(writer, request)
					return
				}

				hahahaHandler("")(writer, request)
			})
		})

		authRoutes(r)
	})
}

func setupSessionStorage() sessions.Store {
	appURL := util.GetStringFromEnvOrFatalError("APP_VK_URL")
	sessionSecret := util.GetStringFromEnvOrFatalError("SESSION_SECRET")

	maxAge := 86400 * 30
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
