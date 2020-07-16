package main

import (
	"compress/flate"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/vk"
	"github.com/crossworth/cartola-web-admin/web"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(flate.DefaultCompression))

	vkClient, err := vk.NewVKClient("", "")
	if err != nil {
		log.Printf("erro ao criar o cliente VK, %v", err)
	}

	router.Mount("/api", api.NewServer(vkClient))
	router.Mount("/", web.New())

	log.Println("iniciando o servidor na porta 8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("erro ao iniciar o servidor http, %v", err)
	}
}
