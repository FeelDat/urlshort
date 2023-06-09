package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func GetFullURL(repository storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := mux.Vars(r)["id"]
		v, err := repository.GetFullURL(shortURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Location", v)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func ShortenURL(repository storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullURL, err := io.ReadAll(r.Body)
		if err != nil || len(fullURL) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		response := "http://" + r.Host + "/" + repository.ShortenURL(string(fullURL))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))

	}
}
