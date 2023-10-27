// Package storage provides functions and types for storing and managing shortened URLs.
//
// This package is used for in-memory storage of URLs and supports batch processing.

package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/utils"
	"go.uber.org/zap"
	"math/rand"
	"os"
)

// URLInfo represents the details of a shortened URL.

type URLInfo struct {
	UUID        string `json:"uuid"`         // The UUID of the user who shortened the URL.
	ShortURL    string `json:"short_url"`    // The shortened version of the URL.
	OriginalURL string `json:"original_url"` // The original URL.
}

// storage is an in-memory storage structure for URLs.
type storage struct {
	Links    map[string]string             // Maps from shortened URLs to original URLs.
	UserURLs map[string][]models.UsersURLS // A map of user IDs to their list of shortened URLs.
	file     *os.File                      // The file for storing the URLs.
}

// NewInMemStorage initializes an in-memory storage for URLs.
//
// If a filePath is provided, it will be used to open a file where the URLs can be saved.
// Returns an error if there's an issue opening the file.
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

// DeleteURLS deletes a list of URLs for a given user.
//
// This function doesn't return anything but can be expanded to return an error if needed.

func (s *storage) DeleteURLS(ctx context.Context, userID string, shortLinks []string, logger *zap.SugaredLogger) {
}

// GetUsersURLS retrieves all the URLs for a given user.
//
// Returns a list of URLs or an error if no URLs are found for the given user.

func (s *storage) GetUsersURLS(ctx context.Context, userID string, baseAddr string) ([]models.UsersURLS, error) {
	if urls, ok := s.UserURLs[userID]; ok {
		return urls, nil
	}

	return nil, errors.New("no URLs found for the given userID")
}

// ShortenURL creates a shortened URL for the given original URL.
//
// Returns the shortened URL or an error if the process fails.

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

// GetFullURL retrieves the original URL for a given shortened URL.
//
// Returns the original URL or an error if the shortened URL does not exist.

func (s *storage) GetFullURL(ctx context.Context, shortLink string) (string, error) {
	val, ok := s.Links[shortLink]
	if !ok {
		return "", errors.New("link does not exist")
	}
	return val, nil
}

// ShortenURLBatch creates shortened URLs for a batch of original URLs.
//
// Returns a list of responses containing the shortened URLs or an error if the process fails.

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
		s.UserURLs[uid.(string)] = append(s.UserURLs[uid.(string)], models.UsersURLS{OriginalURL: req.OriginalURL, ShortURL: urlID})

		responses[i] = models.URLRBatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseAddr + "/" + urlID,
		}
	}

	return responses, nil
}
