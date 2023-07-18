package storage

import (
	"database/sql"
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
	err := d.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
