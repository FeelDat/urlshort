package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func GetFullUrl(repository storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := mux.Vars(r)["id"]
		v, err := repository.GetFullUrl(shortURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Location", v)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func ShortenUrl(repository storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullUrl, err := io.ReadAll(r.Body)
		if err != nil || len(fullUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		response := "http://" + r.Host + "/" + repository.ShortenUrl(string(fullUrl))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))

	}
}
