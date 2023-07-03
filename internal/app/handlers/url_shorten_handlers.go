package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonReply struct {
	Result string `json:"result"`
}

type HandlerInterface interface {
	GetFullURL(w http.ResponseWriter, r *http.Request)
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ShortenURLJSON(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repo        storage.Repository
	baseAddress string
}

func NewHandler(repo storage.Repository, baseAddress string) HandlerInterface {
	return &handler{
		repo:        repo,
		baseAddress: baseAddress,
	}
}

func (h *handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	var request jsonRequest
	var reply jsonReply

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(h.baseAddress, "http://") && !strings.HasPrefix(h.baseAddress, "https://") {
		h.baseAddress = "http://" + h.baseAddress
	}

	reply.Result = h.baseAddress + "/" + h.repo.ShortenURL(string(request.URL))

	resp, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(resp))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h *handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
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
