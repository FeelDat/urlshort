package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"time"
)

type Repository interface {
	ShortenURL(ctx context.Context, fullLink string) (string, error)
	GetFullURL(ctx context.Context, shortLink string) (string, error)
	Ping() error
	Close() error
}

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type storage struct {
	Links   map[string]string
	file    *os.File
	encoder *json.Encoder
	db      *sql.DB
}

func NewStorage(ctx context.Context, filePath string, db *sql.DB) (Repository, error) {

	var file *os.File
	var err error
	//Does it need to be off if "" ?

	if db != nil {
		ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		_, err = db.ExecContext(ctrl, "CREATE TABLE IF NOT EXISTS urls(id serial primary key, uuid varchar(36), short_url varchar(20), original_url text)")
		if err != nil {
			return nil, err
		}
	}

	if filePath != "" {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}

	return &storage{
		Links:   make(map[string]string),
		file:    file,
		encoder: json.NewEncoder(file),
		db:      db,
	}, err

}
func (s *storage) Ping() error {
	if s.db == nil {
		return errors.New("no db connection")
	}
	return s.db.Ping()
}

func (s *storage) ShortenURL(ctx context.Context, fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	uid := uuid.NewString()

	if s.db != nil {
		ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		_, err := s.db.ExecContext(ctrl, `INSERT INTO urls(uuid, short_url, original_url) VALUES($1, $2, $3)`, uid, urlID, fullLink)
		if err != nil {
			return "", err
		}
	} else if s.file != nil {
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
	} else {
		s.Links[urlID] = fullLink
	}

	return urlID, nil
}

func (s *storage) GetFullURL(ctx context.Context, shortLink string) (string, error) {

	if s.db != nil {
		var originalURL string
		ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		err := s.db.QueryRowContext(ctrl, `SELECT original_url FROM urls WHERE short_url = $1`, shortLink).Scan(&originalURL)
		if err != nil {
			if err == sql.ErrNoRows {
				return "", errors.New("link does not exist")
			}
			return "", err
		}

		return originalURL, nil
	} else {
		val, ok := s.Links[shortLink]
		if !ok {
			return "", errors.New("link does not exist")
		}
		return val, nil
	}
}

func (s *storage) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
