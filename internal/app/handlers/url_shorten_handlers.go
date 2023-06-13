package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type HandlerInterface interface {
	GetFullURL(w http.ResponseWriter, r *http.Request)
	ShortenURL(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	repo        storage.Repository
	baseAddress string
}

func NewHandler(repo storage.Repository, baseAddress string) HandlerInterface {
	return &Handler{
		repo:        repo,
		baseAddress: baseAddress,
	}
}

func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v, err := h.repo.GetFullURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	fullURL, err := io.ReadAll(r.Body)
	if err != nil || len(fullURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(h.baseAddress, "http://") && !strings.HasPrefix(h.baseAddress, "https://") {
		h.baseAddress = "http://" + h.baseAddress
	}

	response := h.baseAddress + "/" + h.repo.ShortenURL(string(fullURL))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

//func GetFullURL(repository storage.Repository) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		shortURL := chi.URLParam(r, "id")
//		if shortURL == "" {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		v, err := repository.GetFullURL(shortURL)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		w.Header().Set("Location", v)
//		w.WriteHeader(http.StatusTemporaryRedirect)
//	}
//}

//func ShortenURL(repository storage.Repository, baseAddr string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		fullURL, err := io.ReadAll(r.Body)
//		if err != nil || len(fullURL) == 0 {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		if !strings.HasPrefix(baseAddr, "http://") && !strings.HasPrefix(baseAddr, "https://") {
//			baseAddr = "http://" + baseAddr
//		}
//
//		response := baseAddr + "/" + repository.ShortenURL(string(fullURL))
//
//		w.Header().Set("Content-Type", "text/plain")
//		w.WriteHeader(http.StatusCreated)
//		w.Write([]byte(response))
//
//	}
//}
