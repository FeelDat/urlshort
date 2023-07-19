package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Repository interface {
	ShortenURL(ctx context.Context, fullLink string) (string, error)
	GetFullURL(ctx context.Context, shortLink string) (string, error)
}

type dbStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, db *sql.DB) (Repository, error) {

	ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	_, err := db.ExecContext(ctrl, "CREATE TABLE IF NOT EXISTS urls(id serial primary key, uuid varchar(36), short_url varchar(20), original_url text)")
	if err != nil {
		return nil, err
	}

	return &dbStorage{db: db}, err

}

func (s *dbStorage) ShortenURL(ctx context.Context, fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	uid := uuid.NewString()

	ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	_, err := s.db.ExecContext(ctrl, `INSERT INTO urls(uuid, short_url, original_url) VALUES($1, $2, $3)`, uid, urlID, fullLink)
	if err != nil {
		return "", err
	}

	return urlID, nil
}

func (s *dbStorage) GetFullURL(ctx context.Context, shortLink string) (string, error) {

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
}
