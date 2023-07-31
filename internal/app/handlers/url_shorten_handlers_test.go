package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type mockStorage struct {
	Links   map[string]string
	Users   map[string][]models.UsersURLS
	file    *os.File
	encoder *json.Encoder
}

func (m *mockStorage) GetUsersURLS(ctx context.Context, userID string) ([]models.UsersURLS, error) {
	if urls, ok := m.Users[userID]; ok {
		return urls, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockStorage) ShortenURL(ctx context.Context, fullURL string, userID string) (string, error) {
	shortURL := "UySmre7XjFr"
	m.Links[shortURL] = fullURL
	m.Users[userID] = append(m.Users[userID], models.UsersURLS{ShortURL: shortURL, OriginalURL: fullURL})
	return shortURL, nil
}

func newTestToken() (string, error) {
	key := "8PNHgjK2kPunGpzMgL0ZmMdJCRKy2EnL/Cg0GbnELLI="

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	userID := "testUserID"

	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}

func TestShortenURL(t *testing.T) {
	testCases := []struct {
		name                string
		longLink            string
		method              string
		expectedStatusCode  int
		expectedContentType string
		authenticated       bool
	}{
		{
			name:                "authenticated request",
			longLink:            "https://practicum.yandex.ru/",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: "text/plain",
			authenticated:       true,
		},
		{
			name:                "unauthenticated request",
			longLink:            "https://practicum.yandex.ru/",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: "text/plain; charset=utf-8",
			authenticated:       false,
		},
	}

	mockStorage, _ := storage.NewInMemStorage("short-url-db.json")
	mockHandler := NewHandler(mockStorage, "localhost:8080", nil)

	token, err := newTestToken()
	require.NoError(t, err)

	defer os.Remove("short-url-db.json")

	router := chi.NewRouter()
	router.Post("/", mockHandler.ShortenURL)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, ts.URL+"/", strings.NewReader(tt.longLink))
			require.NoError(t, err)

			if tt.authenticated {
				r.AddCookie(&http.Cookie{
					Name:  "jwt",
					Value: token,
				})
			}

			resp, err := ts.Client().Do(r)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))

			if tt.authenticated {
				urls, err := mockStorage.GetUsersURLS(context.Background(), "testUserID")
				require.NoError(t, err)
				assert.Len(t, urls, 1)
				assert.Equal(t, urls[0].OriginalURL, tt.longLink)
			}
		})
	}
}
