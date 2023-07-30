package storage

import (
	"context"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/utils"
	"math/rand"
)

type MemoryStorage struct {
	urlMap  map[string]*models.URL
	userMap map[string][]*models.URL
}

func NewMemoryStorage() Repository {
	return &MemoryStorage{
		urlMap:  make(map[string]*models.URL),
		userMap: make(map[string][]*models.URL),
	}
}

func (m *MemoryStorage) GetUsersURLS(ctx context.Context, userID string) ([]models.UsersURLS, error) {
	userUrls, exists := m.userMap[userID]
	if !exists {
		return nil, errors.New("no urls found for the user")
	}

	var urls []models.UsersURLS

	for _, url := range userUrls {
		urls = append(urls, models.UsersURLS{
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}

	return urls, nil
}

func (m *MemoryStorage) ShortenURL(ctx context.Context, fullLink string) (string, error) {
	shortURL := utils.Base62Encode(rand.Uint64())
	userID := ctx.Value("userID").(string)

	if _, exists := m.urlMap[shortURL]; exists {
		return "", errors.New("short url already exists")
	}

	url := &models.URL{
		ShortURL:    shortURL,
		OriginalURL: fullLink,
		UserID:      userID,
	}

	m.urlMap[shortURL] = url
	m.userMap[userID] = append(m.userMap[userID], url)

	return shortURL, nil
}

func (m *MemoryStorage) GetFullURL(ctx context.Context, shortLink string) (string, error) {
	url, exists := m.urlMap[shortLink]
	if !exists {
		return "", errors.New("short url does not exist")
	}

	return url.OriginalURL, nil
}

func (m *MemoryStorage) ShortenURLBatch(ctx context.Context, batch []models.URLBatchRequest, baseAddr string) ([]models.URLRBatchResponse, error) {
	if len(batch) == 0 {
		return nil, errors.New("empty batch")
	}

	userID := ctx.Value("userID").(string)
	responses := make([]models.URLRBatchResponse, len(batch))

	for i, req := range batch {
		shortURL := utils.Base62Encode(rand.Uint64())

		if _, exists := m.urlMap[shortURL]; exists {
			return nil, errors.New("short url already exists")
		}

		url := &models.URL{
			ShortURL:    shortURL,
			OriginalURL: req.OriginalURL,
			UserID:      userID,
		}

		m.urlMap[shortURL] = url
		m.userMap[userID] = append(m.userMap[userID], url)

		responses[i] = models.URLRBatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseAddr + "/" + shortURL,
		}
	}

	return responses, nil
}
