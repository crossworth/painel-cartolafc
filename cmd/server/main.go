package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/cartola"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/updater"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
	"github.com/crossworth/cartola-web-admin/vk/openid"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("erro carregando arquivo .env, %v", err)
	}

	logger.Setup(logger.LogInfo, "cartola_web_admin.log")

	logger.Log.Info().Msg("iniciando servidor")
	vkClient := setupVKClient()
	db := setupDatabase()
	session := setupSessionStorage()
	appCache := cache.NewCache()
	topicUpdater := updater.NewTopicUpdater(db.GetDB())

	appName := util.GetStringFromEnvOrFatalError("APP_NAME")

	vkBotQuoteID := util.GetIntFromEnvOrFatalError("VK_BOT_QUOTES_ID")

	vkAppID := util.GetStringFromEnvOrFatalError("VK_APP_ID")
	vkSecureKey := util.GetStringFromEnvOrFatalError("VK_SECURE_KEY")
	vkCallBackURL := util.GetStringFromEnvOrFatalError("VK_CALLBACK_URL")

	goth.UseProviders(
		openid.New(vkAppID, vkSecureKey, vkCallBackURL, "groups"),
	)

	gothic.Store = session

	superAdminsStr := util.GetStringFromEnvOrFatalError("SUPER_ADMIN_VK_IDS")
	superAdmins := util.StringToIntSlice(superAdminsStr)

	if len(superAdmins) == 0 {
		logger.Log.Warn().Msg("nenhum super administrador definido")
	} else {
		logger.Log.Info().Interface("super_admins", superAdmins).Msg("super administradores definidos")
	}

	app := cartola.NewCartola(appName, vkClient, db, session, appCache, topicUpdater, superAdmins, vkBotQuoteID)

	// topicUpdater.RegisterWorker(func(job updater.TopicUpdateJob) error {
	// 	logger.Log.Info().Msgf("Handling1 job: %d", job.ID)
	// 	time.Sleep(10 * time.Millisecond)
	// 	return fmt.Errorf("aaa isso é um erro")
	// }, true)
	// topicUpdater.StartProcessing()

	appPort := util.GetStringFromEnvOrFatalError("APP_PORT")
	logger.Log.Info().Msgf("iniciando o servidor na porta %s", appPort)

	err = http.ListenAndServe(fmt.Sprintf(":%s", appPort), app)
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao iniciar o servidor http, %v", err)
	}
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
