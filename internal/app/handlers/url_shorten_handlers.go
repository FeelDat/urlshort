package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
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
	Ping(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	inMemoryRepo   storage.InMemoryRepository
	inDatabaseRepo storage.DatabaseRepository
	baseAddress    string
}

func NewHandler(inMemoryRepo storage.InMemoryRepository, baseAddress string, inDatabase storage.DatabaseRepository) HandlerInterface {
	return &handler{
		inMemoryRepo:   inMemoryRepo,
		inDatabaseRepo: inDatabase,
		baseAddress:    baseAddress,
	}
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {

	err := h.inDatabaseRepo.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v, err := h.inMemoryRepo.GetFullURL(shortURL)
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

	h.baseAddress = utils.AddPrefix(h.baseAddress)

	shortURL, err := h.inMemoryRepo.ShortenURL(string(request.URL))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	reply.Result = h.baseAddress + "/" + shortURL

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

	h.baseAddress = utils.AddPrefix(h.baseAddress)

	shortURL, err := h.inMemoryRepo.ShortenURL(string(fullURL))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	response := h.baseAddress + "/" + shortURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
