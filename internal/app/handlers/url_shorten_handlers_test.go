package handlers

import (
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/config"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type mockStorage struct {
	Links   map[string]string
	file    *os.File
	encoder *json.Encoder
}

func newMockMemoryStorage() (storage.Repository, error) {
	return &mockStorage{
		Links: make(map[string]string),
	}, nil
}

func (m *mockStorage) ShortenURL(fullURL string) (string, error) {

	shortURL := "UySmre7XjFr"
	m.Links[shortURL] = fullURL
	return shortURL, nil
}

func (m *mockStorage) GetFullURL(shortLink string) (string, error) {
	val, ok := m.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}

func (m *mockStorage) Close() error {
	if m.file != nil {
		return m.file.Close()
	}
	return nil
}

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

	mckStorage, _ := newMockMemoryStorage()
	_, err := mckStorage.ShortenURL("https://practicum.yandex.ru/")
	require.NoError(t, err)

	mockHandler := NewHandler(mckStorage, &config.Config{})

	router := chi.NewRouter()
	router.Get("/{id}", mockHandler.GetFullURL)

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

	mockStorage, _ := storage.NewInMemoryStorage("short-url-db.json")
	mockHandler := NewHandler(mockStorage, &config.Config{
		ServerAddress: "localhost:8080",
		BaseAddress:   "localhost:8080",
		FilePath:      ""})

	defer os.Remove("short-url-db.json")

	router := chi.NewRouter()
	router.Post("/", mockHandler.ShortenURL)

	ts := httptest.NewServer(router)
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

func Test_handler_ShortenURLJSON(t *testing.T) {

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

	mockStorage, _ := storage.NewInMemoryStorage("short-url-db.json")
	mockHandler := NewHandler(mockStorage, &config.Config{
		ServerAddress: "localhost:8080",
		BaseAddress:   "localhost:8080",
		FilePath:      ""})

	defer os.Remove("short-url-db.json")

	router := chi.NewRouter()
	router.Post("/", mockHandler.ShortenURL)

	ts := httptest.NewServer(router)
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
