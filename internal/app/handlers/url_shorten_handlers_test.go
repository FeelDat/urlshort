package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFullURL(t *testing.T) {
	testCases := []struct {
		name               string
		shortLink          string
		method             string
		expectedStatusCode int
		expectedLink       string
	}{
		{
			name:               "successful request",
			shortLink:          "UySmre7XjFr",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLink:       "https://practicum.yandex.ru/",
		},
		{
			name:               "not existing link",
			shortLink:          "Usuhf784",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedLink:       "",
		},
		{
			name:               "empty id",
			shortLink:          "",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedLink:       "",
		},
	}

	mockStorage := storage.InitInMemoryStorage()
	mockStorage.Links["UySmre7XjFr"] = "https://practicum.yandex.ru/"

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tt.shortLink})
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetFullUrl(mockStorage))
			h(w, r)

			result := w.Result()
			assert.Equal(t, tt.expectedLink, result.Header.Get("Location"))
			assert.Equal(t, tt.expectedStatusCode, result.StatusCode)
		})
	}

}

func TestShortenURL(t *testing.T) {
	testCases := []struct {
		name                string
		longLink            string
		method              string
		expectedStatusCode  int
		expectedContentType string
	}{
		{
			name:                "successful test",
			longLink:            "https://practicum.yandex.ru/",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: "text/plain",
		},
		{
			name:                "no link",
			longLink:            "",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: "",
		},
	}

	mockStorage := storage.InitInMemoryStorage()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.longLink))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenUrl(mockStorage))
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.expectedStatusCode, result.StatusCode)
			assert.Equal(t, tt.expectedContentType, result.Header.Get("Content-Type"))

		})
	}
}
