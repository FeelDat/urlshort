package storage

import (
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"math/rand"
)

type Repository interface {
	ShortenUrl(fullLink string) string
	GetFullUrl(shortLink string) (string, error)
}
type InMemoryStorage struct {
	Links map[string]string
}

func InitInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Links: make(map[string]string),
	}
}

func (mapStorage InMemoryStorage) ShortenUrl(fullLink string) string {

	urlID := utils.Base62Encode(rand.Uint64())
	mapStorage.Links[urlID] = string(fullLink)
	return urlID
}

func (mapStorage InMemoryStorage) GetFullUrl(shortLink string) (string, error) {

	val, ok := mapStorage.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}
