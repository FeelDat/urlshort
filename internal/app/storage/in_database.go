package storage

import (
	"context"
	"database/sql"
	"time"
)

type DatabaseStorage struct {
	db *sql.DB
}

type DatabaseRepository interface {
	Ping() error
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{
		db: db,
	}
}

func (d *DatabaseStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := d.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
