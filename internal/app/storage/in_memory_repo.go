package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/utils"
	"math/rand"
	"os"
)

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type storage struct {
	Links    map[string]string
	UserURLs map[string][]models.UsersURLS
	file     *os.File
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
		Links:    make(map[string]string),
		UserURLs: make(map[string][]models.UsersURLS),
		file:     file,
	}, err
}

func (s *storage) GetUsersURLS(ctx context.Context, userID string, baseAddr string) ([]models.UsersURLS, error) {
	if urls, ok := s.UserURLs[userID]; ok {
		return urls, nil
	}

	return nil, errors.New("no URLs found for the given userID")
}

func (s *storage) ShortenURL(ctx context.Context, fullLink string) (string, error) {
	urlID := utils.Base62Encode(rand.Uint64())
	uid := ctx.Value(models.CtxKey("userID"))

	urlInfo := URLInfo{
		UUID:        uid.(string),
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

	s.UserURLs[uid.(string)] = append(s.UserURLs[uid.(string)], models.UsersURLS{OriginalURL: fullLink, ShortURL: urlID})

	return urlID, nil
}

func (s *storage) GetFullURL(ctx context.Context, shortLink string) (string, error) {
	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}

func (s *storage) ShortenURLBatch(ctx context.Context, batch []models.URLBatchRequest, baseAddr string) ([]models.URLRBatchResponse, error) {
	if len(batch) == 0 {
		return nil, errors.New("empty batch")
	}

	responses := make([]models.URLRBatchResponse, len(batch))
	uid := ctx.Value(models.CtxKey("userID"))

	for i, req := range batch {
		urlID := utils.Base62Encode(rand.Uint64())
		urlInfo := URLInfo{
			UUID:        uid.(string),
			ShortURL:    urlID,
			OriginalURL: req.OriginalURL,
		}

		data, err := json.Marshal(&urlInfo)
		if err != nil {
			return nil, err
		}

		_, err = s.file.Write(data)
		if err != nil {
			return nil, err
		}

		s.Links[urlID] = req.OriginalURL
		// Add the URL to UserURLs
		s.UserURLs[uid.(string)] = append(s.UserURLs[uid.(string)], models.UsersURLS{OriginalURL: req.OriginalURL, ShortURL: urlID})

		responses[i] = models.URLRBatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseAddr + "/" + urlID,
		}
	}

	return responses, nil
}
