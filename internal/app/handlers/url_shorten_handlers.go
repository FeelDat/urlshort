package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type HandlerInterface interface {
	GetFullURL(w http.ResponseWriter, r *http.Request)
	ShortenURL(w http.ResponseWriter, r *http.Request)
	ShortenURLJSON(w http.ResponseWriter, r *http.Request)
	ShortenURLBatch(w http.ResponseWriter, r *http.Request)
	GetUsersURLS(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repository  storage.Repository
	baseAddress string
	logger      *zap.SugaredLogger
}

func NewHandler(repo storage.Repository, baseAddress string, logger *zap.SugaredLogger) HandlerInterface {
	return &handler{
		repository:  repo,
		baseAddress: baseAddress,
		logger:      logger,
	}
}

func (h *handler) GetUsersURLS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	urls, err := h.repository.GetUsersURLS(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		h.logger.Errorw("Failed to unmarshal request body", "error", err)
		return
	}

	h.baseAddress, err = utils.AddPrefix(h.baseAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to add prefix to baseAddress", "error", err)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	cntx := context.WithValue(r.Context(), "userID", userID)

	shortURL, err := h.repository.ShortenURL(cntx, string(request.URL))
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			reply.Result = h.baseAddress + "/" + shortURL
			resp, err := json.Marshal(reply)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.logger.Errorw("Failed to marshal response", "error", err)
				return
			}
			_, err = w.Write([]byte(resp))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.logger.Errorw("Failed to write response", "error", err)
				return
			}
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Errorw("Failed to store shortened URL in DB", "error", err)

			return
		}
	}

	reply.Result = h.baseAddress + "/" + shortURL

	resp, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to marshal response", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(resp))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to write response", "error", err)
		return
	}

}

func (h *handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	fullURL, err := io.ReadAll(r.Body)
	if err != nil || len(fullURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("No URL in request body to shorten", "error", err)
		return
	}

	h.baseAddress, err = utils.AddPrefix(h.baseAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to add prefix to baseAddress", "error", err)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	cntx := context.WithValue(r.Context(), "userID", userID)

	shortURL, err := h.repository.ShortenURL(cntx, string(fullURL))
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			w.WriteHeader(http.StatusConflict)
			response := h.baseAddress + "/" + shortURL
			_, err := w.Write([]byte(response))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.logger.Errorw("Failed to write response", "error", err)
				return
			}
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Errorw("Failed to store shortened URL in DB", "error", err)
			return
		}
	}

	response := h.baseAddress + "/" + shortURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to write response", "error", err)
		return
	}
}

func (h *handler) ShortenURLBatch(w http.ResponseWriter, r *http.Request) {

	urls := make([]models.URLBatchRequest, 0)
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorw("Failed to get URLs from json batch", "error", err)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "empty batch", http.StatusBadRequest)
		h.logger.Errorw("URLs batch is empty", "error", err)
		return
	}

	h.baseAddress, err = utils.AddPrefix(h.baseAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to add prefix to baseAddress", "error", err)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID := cookie.Value

	cntx := context.WithValue(r.Context(), "userID", userID)

	result, err := h.repository.ShortenURLBatch(cntx, urls, h.baseAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorw("Failed to store shortened URLs batch in DB", "error", err)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to marshal response", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorw("Failed to write response", "error", err)
		return
	}

}
