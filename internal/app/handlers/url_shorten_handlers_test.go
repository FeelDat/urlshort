package handlers

import (
	"context"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mocks.NewMockRepository(ctrl)

	// Create a handler with the mock repository
	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())

	testCases := []struct {
		name                string
		longLink            string
		method              string
		expectedStatusCode  int
		expectedContentType string
		authenticated       bool
	}{
		{
			name:                "authenticated request with valid URL",
			longLink:            "https://practicum.yandex.ru/",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: "text/plain",
			authenticated:       true,
		},
		{
			name:                "unauthenticated request with valid URL",
			longLink:            "https://practicum.yandex.ru/",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: "text/plain; charset=utf-8",
			authenticated:       false,
		},
		{
			name:                "authenticated request with invalid URL",
			longLink:            "",
			method:              http.MethodPost,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: "",
			authenticated:       true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := strings.NewReader(tt.longLink)
			req, err := http.NewRequest(tt.method, "/", reqBody)
			require.NoError(t, err)

			if tt.authenticated {
				token, err := newTestToken()
				require.NoError(t, err)
				ctx := context.WithValue(req.Context(), models.CtxKey("userID"), "testUserID")
				req = req.WithContext(ctx)
				req.AddCookie(&http.Cookie{
					Name:  "jwt",
					Value: token,
				})
			}

			if tt.expectedStatusCode == http.StatusCreated {
				// If the request is expected to succeed, set up the expectation
				mockRepo.EXPECT().ShortenURL(gomock.Any(), tt.longLink).Return("shortened", nil)
			} else {
				// If the request is expected to fail, don't expect ShortenURL to be called
				mockRepo.EXPECT().ShortenURL(gomock.Any(), gomock.Any()).Times(0)
			}

			rr := httptest.NewRecorder()
			handler.ShortenURL(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			assert.Equal(t, tt.expectedContentType, rr.Header().Get("Content-Type"))

		})
	}
}

//TODO write test for other handlers

func TestHandlerGetFullURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mocks.NewMockRepository(ctrl)

	// Create a handler with the mock repository
	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())

	testCases := []struct {
		name               string
		shortURL           string
		expectedStatusCode int
		expectedLocation   string
		repoResult         string // Mocked repository result
		repoError          error  // Mocked repository error
	}{
		{
			name:               "Successful Retrieval",
			shortURL:           "localhost:8080/shortened",
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLocation:   "https://practicum.yandex.ru/",
			repoResult:         "https://practicum.yandex.ru/",
		},
		{
			name:               "Link does not exist",
			shortURL:           "localhost:8080/nonexistent",
			expectedStatusCode: http.StatusNotFound,
			expectedLocation:   "",
			repoError:          errors.New("link does not exist"), // Empty result to simulate not found
		},
		{
			name:               "Empty Short URL",
			shortURL:           "",
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "",
			repoResult:         "", // Irrelevant as the URL is empty
		},
		{
			name:               "Link is deleted",
			shortURL:           "localhost:8080/error",
			expectedStatusCode: http.StatusGone,
			expectedLocation:   "",
			repoError:          errors.New("link is deleted"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the repository behavior
			if tt.shortURL != "" {
				mockRepo.EXPECT().GetFullURL(gomock.Any(), tt.shortURL).Return(tt.repoResult, tt.repoError)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			// Set the URL parameter using chi.URLParam
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tt.shortURL)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			handler.GetFullURL(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedStatusCode == http.StatusFound {
				assert.Equal(t, tt.expectedLocation, rr.Header().Get("Location"))
			}
		})
	}
}

//
//func TestHandlerShortenURLJSON(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// Create a mock repository
//	mockRepo := mocks.NewMockRepository(ctrl)
//
//	// Create a handler with the mock repository
//	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())
//
//	testCases := []struct {
//		name               string
//		longLink           string
//		expectedStatusCode int
//		expectedJSON       string
//		repoResult         string // Mocked repository result
//		repoError          error  // Mocked repository error
//	}{
//		{
//
//		}
//	}
//
//}

//
//func TestHandlerShortenURLBatch(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// Create a mock repository
//	mockRepo := mocks.NewMockRepository(ctrl)
//
//	// Create a handler with the mock repository
//	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())
//
//}
//
//func TestHandlerGetUsersURLS(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// Create a mock repository
//	mockRepo := mocks.NewMockRepository(ctrl)
//
//	// Create a handler with the mock repository
//	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())
//
//}
//
//func TestHandlerDeleteURLS(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// Create a mock repository
//	mockRepo := mocks.NewMockRepository(ctrl)
//
//	// Create a handler with the mock repository
//	handler := NewHandler(mockRepo, "localhost:8080", zap.NewNop().Sugar())
//
//}
