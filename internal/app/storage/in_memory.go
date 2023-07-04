package storage

import (
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/google/uuid"
	"math/rand"
	"os"
)

type Repository interface {
	ShortenURL(fullLink string) (string, error)
	GetFullURL(shortLink string) (string, error)
}

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type inMemoryStorage struct {
	Links    map[string]string
	filePath string
}

func NewInMemoryStorage(filePath string) Repository {
	return &inMemoryStorage{
		Links:    make(map[string]string),
		filePath: filePath,
	}
}

func (s *inMemoryStorage) ShortenURL(fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	s.Links[urlID] = string(fullLink)

	if s.filePath != "" {
		f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return "", err
		}
		defer f.Close()

		urlInfo := URLInfo{
			UUID:        uuid.NewString(),
			ShortURL:    urlID,
			OriginalURL: fullLink,
		}

		data, err := json.Marshal(&urlInfo)
		if err != nil {
			return "", err
		}

		f.Write(data)
	}

	return urlID, nil
}

func (s *inMemoryStorage) GetFullURL(shortLink string) (string, error) {

	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}
