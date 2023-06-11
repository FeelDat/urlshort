package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func GetFullURL(repository storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "id")
		if shortURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v, err := repository.GetFullURL(shortURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", v)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func ShortenURL(repository storage.Repository, baseAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullURL, err := io.ReadAll(r.Body)
		if err != nil || len(fullURL) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response := "http://" + baseAddr + "/" + repository.ShortenURL(string(fullURL))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))

	}
}
