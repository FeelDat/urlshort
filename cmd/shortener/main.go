package main

import (
	"github.com/FeelDat/urlshort/internal/app/handlers"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	mapStorage := storage.InitInMemoryStorage()

	router.HandleFunc(`/`, handlers.ShortenURL(mapStorage)).Methods("POST")
	router.HandleFunc(`/{id}`, handlers.GetFullURL(mapStorage)).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
