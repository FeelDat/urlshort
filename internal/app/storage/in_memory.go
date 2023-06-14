package storage

import (
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"math/rand"
)

type Repository interface {
	ShortenURL(fullLink string) string
	GetFullURL(shortLink string) (string, error)
}

type inMemoryStorage struct {
	Links map[string]string
}

func NewInMemoryStorage() Repository {
	return &inMemoryStorage{
		Links: make(map[string]string),
	}
}

func (s *inMemoryStorage) ShortenURL(fullLink string) string {

	urlID := utils.Base62Encode(rand.Uint64())
	s.Links[urlID] = string(fullLink)
	return urlID
}

func (s *inMemoryStorage) GetFullURL(shortLink string) (string, error) {

	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}
