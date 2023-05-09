package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-gorp/gorp/v3"
	"github.com/heppu/miniurl/storage"
	"github.com/lib/pq"
)

type record struct {
	Hash string `db:"hash, primarykey"`
	Url  string `db:"url"`
}

type Storage struct {
	db *gorp.DbMap
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	dbmap.AddTableWithName(record{}, "records")
	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create missing tables: %w", err)
	}

	return &Storage{db: dbmap}, nil
}

func (s *Storage) AddUrl(url, hash string) error {
	err := s.db.Insert(&record{Hash: hash, Url: url})
	if !violatesUniqueConstrain(err) {
		return err
	}

	oldUrl, err := s.GetUrl(hash)
	if err != nil {
		return err
	}

	if oldUrl != url {
		return storage.ErrHashCollision
	}

	return nil
}

func (s *Storage) GetUrl(hash string) (string, error) {
	url, err := s.db.SelectStr("SELECT url FROM records WHERE hash = $1", hash)
	if err != nil {
		return "", err
	}
	if url == "" {
		return "", storage.ErrUrlNotFound
	}

	return url, nil
}

func violatesUniqueConstrain(err error) bool {
	const uniqErrCode = "23505"
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && string(pgErr.Code) == uniqErrCode
}
