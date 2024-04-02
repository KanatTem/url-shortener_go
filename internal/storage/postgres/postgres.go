package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("postgres", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.NewStorage.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO urls (alias, url) VALUES ($1, $2)`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		var pqErr pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {

			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)

	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: get last insert ID: %w", op, err)
	}

	return id, nil

}

func (s *Storage) getUrl(alias string) (string, error) {
	const op = "storage.postgres.NewStorage.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM urls WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var url string

	err = stmt.QueryRow(alias).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return url, nil

}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.postgres.NewStorage.DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM urls WHERE alias = $1`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrURLNotFound
	}
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
