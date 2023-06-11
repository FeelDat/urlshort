package handlers

import (
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			expectedStatusCode: http.StatusNotFound,
			expectedLink:       "",
		},
	}

	mockStorage := storage.InitInMemoryStorage()
	mockStorage.Links["UySmre7XjFr"] = "https://practicum.yandex.ru/"

	router := chi.NewRouter()
	router.Get("/{id}", GetFullURL(mockStorage))

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			r, err := http.NewRequest(tt.method, ts.URL+"/"+tt.shortLink, nil)
			require.NoError(t, err)

			client := &http.Client{
				// Prevent auto-following of redirects
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(r)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedLink, resp.Header.Get("Location"))
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
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
	mockAddress := "localhost:8888"

	ts := httptest.NewServer(ShortenURL(mockStorage, mockAddress))
	defer ts.Close()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			r, err := http.NewRequest(tt.method, ts.URL+"/", strings.NewReader(tt.longLink))
			require.NoError(t, err)

			resp, err := ts.Client().Do(r)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))

		})
	}
}
