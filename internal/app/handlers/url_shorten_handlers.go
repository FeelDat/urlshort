package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
)

type HandlerInterface interface {
	GetFullURL(w http.ResponseWriter, r *http.Request)
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ShortenURLJSON(w http.ResponseWriter, r *http.Request)
	ShortenURLBatch(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repository  storage.Repository
	baseAddress string
}

func NewHandler(repo storage.Repository, baseAddress string) HandlerInterface {
	return &handler{
		repository:  repo,
		baseAddress: baseAddress,
	}
}

func (h *handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v, err := h.repository.GetFullURL(r.Context(), shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *handler) ShortenURLJSON(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	var request models.JSONRequest
	var reply models.JSONResponse

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

	shortURL, err := h.repository.ShortenURL(r.Context(), string(request.URL))
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			reply.Result = h.baseAddress + "/" + shortURL
			resp, err := json.Marshal(reply)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_, err = w.Write([]byte(resp))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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

	shortURL, err := h.repository.ShortenURL(r.Context(), string(fullURL))
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			w.WriteHeader(http.StatusConflict)
			response := h.baseAddress + "/" + shortURL
			_, err := w.Write([]byte(response))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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

func (h *handler) ShortenURLBatch(w http.ResponseWriter, r *http.Request) {

	urls := make([]models.URLBatchRequest, 0)
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "empty batch", http.StatusBadRequest)
		return
	}

	h.baseAddress = utils.AddPrefix(h.baseAddress)

	result, err := h.repository.ShortenURLBatch(r.Context(), urls, h.baseAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
