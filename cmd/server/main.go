package main

import (
	"compress/flate"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/crossworth/cartola-web-admin/api"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/typesense"
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
	router.Use(middleware.BasicAuth("CartolaFC", creds))

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

	go func() {
		return
		ts, err := typesense.NewSearch("192.168.0.62", "8108", "teste")
		if err != nil {
			log.Fatalln(err)
		}

		err = ts.CreateCollections()
		if err != nil {
			log.Fatalln(err)
		}

		profiles, err := db.ProfilesAll(context.TODO())
		if err != nil {
			log.Fatalln(err)
		}

		lenProfiles := len(profiles)

		for i, profile := range profiles {
			if i%100 == 0 {
				log.Printf("inserindo perfil %d/%d\n", i, lenProfiles)
			}
			err = ts.InsertProfile(typesense.ToTypeSenseProfile(profile))
			if err != nil {
				log.Printf("erro ao inserir perfil %v\n", err)
			}
		}

		log.Printf("inserido %d perfils\n", lenProfiles)

		topics, err := db.TopicsAll(context.TODO())
		if err != nil {
			log.Fatalln(err)
		}

		lenTopics := len(topics)

		for i, topic := range topics {
			if i%100 == 0 {
				log.Printf("inserindo tópico %d/%d\n", i, lenTopics)
			}
			err = ts.InsertTopic(typesense.ToTypeSenseTopic(topic))
			if err != nil {
				log.Printf("erro ao inserir tópico %v\n", err)
			}
		}

		log.Printf("inserido %d tópicos\n", lenTopics)

		comments, err := db.CommentsAll(context.TODO())
		if err != nil {
			log.Fatalln(err)
		}

		lenComments := len(comments)

		for i, comment := range comments {
			if i%100 == 0 {
				log.Printf("inserindo comentário %d/%d\n", i, lenComments)
			}
			err = ts.InsertComment(typesense.ToTypeSenseComment(comment))
			if err != nil {
				log.Printf("erro ao inserir comentário %v\n", err)
			}
		}

		log.Printf("inserido %d comentários\n", lenComments)

	}()

	router.Mount("/api", api.NewServer(vkClient, db))
	router.Mount("/", web.New())

	log.Printf("iniciando o servidor na porta %s\n", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router)
	if err != nil {
		log.Fatalf("erro ao iniciar o servidor http, %v", err)
	}
}
