package main

import (
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	r := chi.NewRouter()
	mapStorage := storage.InitInMemoryStorage()

	baseAddr := conf.BaseAddress
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.ShortenURL(mapStorage, baseAddr))
		r.Get("/{id}", handlers.GetFullURL(mapStorage))
	})

	err = http.ListenAndServe(conf.ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
