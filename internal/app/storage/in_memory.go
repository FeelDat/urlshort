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
	Close() error
}

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type inMemoryStorage struct {
	Links   map[string]string
	file    *os.File
	encoder *json.Encoder
}

func NewInMemoryStorage(filePath string) (Repository, error) {

	var file *os.File
	var err error
	//Does it need to be off if "" ?
	if filePath != "" {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}

	return &inMemoryStorage{
		Links:   make(map[string]string),
		file:    file,
		encoder: json.NewEncoder(file),
	}, err

}

func (s *inMemoryStorage) ShortenURL(fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	s.Links[urlID] = string(fullLink)

	urlInfo := URLInfo{
		UUID:        uuid.NewString(),
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

	return urlID, nil
}

func (s *inMemoryStorage) GetFullURL(shortLink string) (string, error) {

	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}

func (s *inMemoryStorage) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
