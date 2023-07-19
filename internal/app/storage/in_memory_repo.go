package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/google/uuid"
	"math/rand"
	"os"
)

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type storage struct {
	Links map[string]string
	file  *os.File
}

func NewInMemStorage(filePath string) (Repository, error) {

	var file *os.File
	var err error

	if filePath != "" {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}

	return &storage{
		Links: make(map[string]string),
		file:  file,
	}, err

}

func (s *storage) ShortenURL(ctx context.Context, fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	uid := uuid.NewString()

	urlInfo := URLInfo{
		UUID:        uid,
		ShortURL:    urlID,
		OriginalURL: fullLink,
	}

	data, err := json.Marshal(&urlInfo)
	if err != nil {
		return "", err
	}

	_, err = s.file.Write(data)
	if err != nil {
		return "", err
	}
	s.Links[urlID] = fullLink

	return urlID, nil
}

func (s *storage) GetFullURL(ctx context.Context, shortLink string) (string, error) {

	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil

}