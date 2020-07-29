package main

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/vk"
	"github.com/crossworth/cartola-web-admin/web"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("erro carregando arquivo .env, %v", err)
	}

	log.Println("iniciando servidor")

	creds := make(map[string]string)
	creds["admin"] = "cartolavk"

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(flate.DefaultCompression))
	// router.Use(middleware.BasicAuth("CartolaFC", creds))

	log.Println("criando cliente VK")
	vkClient, err := vk.NewVKClient(os.Getenv("VK_EMAIL"), os.Getenv("VK_PASSWORD"))
	if err != nil {
		log.Fatalf("erro ao criar o cliente VK, %v", err)
	}

	log.Println("montando endpoints")

	db, err := database.NewPostgreSQL(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		log.Fatalf("erro ao conectar ao banco de dados, %v", err)
	}

	router.Mount("/api", api.NewServer(vkClient, db))
	router.Mount("/", web.New())

	log.Printf("iniciando o servidor na porta %s\n", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router)
	if err != nil {
		log.Fatalf("erro ao iniciar o servidor http, %v", err)
	}
}
