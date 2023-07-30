package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/FeelDat/urlshort/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"math/rand"
	"time"
)

type Repository interface {
	ShortenURL(ctx context.Context, fullLink string) (string, error)
	GetFullURL(ctx context.Context, shortLink string) (string, error)
	ShortenURLBatch(ctx context.Context, batch []models.URLBatchRequest, baseAddr string) ([]models.URLRBatchResponse, error)
	GetUsersURLS(ctx context.Context, userID string) ([]models.UsersURLS, error)
}

type dbStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) Repository {
	return &dbStorage{db: db}
}

func InitDB(ctx context.Context, db *sql.DB) error {
	ctrl, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = db.ExecContext(ctrl, "CREATE TABLE IF NOT EXISTS urls(id serial primary key, uuid varchar(36), short_url varchar(20), original_url text)")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctrl, "CREATE UNIQUE INDEX IF NOT EXISTS original_url_unique ON urls(original_url)")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *dbStorage) GetUsersURLS(ctx context.Context, userID string) ([]models.UsersURLS, error) {

	ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := s.db.QueryContext(ctrl, `SELECT short_url, original_url FROM urls WHERE uuid = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.UsersURLS

	for rows.Next() {
		var u models.UsersURLS
		if err := rows.Scan(&u.ShortURL, &u.OriginalURL); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return urls, err
}

func (s *dbStorage) ShortenURL(ctx context.Context, fullLink string) (string, error) {

	urlID := utils.Base62Encode(rand.Uint64())
	uid := ctx.Value("userID")

	ctrl, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	_, err := s.db.ExecContext(ctrl, `INSERT INTO urls(uuid, short_url, original_url) VALUES($1, $2, $3)`, uid, urlID, fullLink)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			var shortURL string
			s.db.QueryRowContext(ctrl, `SELECT short_url FROM urls WHERE original_url = $1`, fullLink).Scan(&shortURL)
			return shortURL, err
		}
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

func (s *dbStorage) ShortenURLBatch(ctx context.Context, batch []models.URLBatchRequest, baseAddr string) ([]models.URLRBatchResponse, error) {

	if len(batch) == 0 {
		return nil, errors.New("empty batch")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	responses := make([]models.URLRBatchResponse, len(batch))

	for i, req := range batch {
		urlID := utils.Base62Encode(rand.Uint64())
		uid := ctx.Value("userID")
		_, err = tx.ExecContext(ctx, `INSERT INTO urls(uuid, short_url, original_url) VALUES($1, $2, $3)`, uid, urlID, req.OriginalURL)
		if err != nil {
			return nil, err
		}

		responses[i] = models.URLRBatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseAddr + "/" + urlID,
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return responses, nil

}
