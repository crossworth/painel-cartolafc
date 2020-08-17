package main

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/vk"
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

	creds := make(map[string]string)
	creds["admin"] = "cartolavk"

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(flate.DefaultCompression))
	router.Use(cors.New(corsOpts).Handler)
	router.Use(middleware.Timeout(10 * time.Minute))
	// router.Use(middleware.BasicAuth("CartolaFC", creds))

	logger.Log.Info().Msg("criando cliente VK")
	vkClient, err := vk.NewVKClient(os.Getenv("VK_EMAIL"), os.Getenv("VK_PASSWORD"))
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao criar o cliente VK, %v", err)
	}

	logger.Log.Info().Msg("conectando no banco de dados")
	db, err := database.NewPostgreSQL(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao conectar ao banco de dados, %v", err)
	}

	logger.Log.Info().Msg("montando endpoints")
	router.Mount("/api", api.NewServer(vkClient, db))
	router.Mount("/", web.New())

	logger.Log.Info().Msgf("iniciando o servidor na porta %s\n", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router)
	if err != nil {
		logger.Log.Fatal().Msgf("erro ao iniciar o servidor http, %v", err)
	}
}
