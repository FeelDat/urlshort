package main

import (
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	mapStorage := storage.InitInMemoryStorage()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.ShortenURL(mapStorage))
		r.Get("/{id}", handlers.GetFullURL(mapStorage))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
