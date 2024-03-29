package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

const (
	maxRetries = 5
	retryDelay = 2 * time.Second
)

func New(connString string) (*Storage, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("postgres", connString)
		if err != nil {
			log.Printf("Connection attempt %d failed: %v", attempt, err)
			time.Sleep(retryDelay)
			continue
		}

		if err = db.Ping(); err == nil {
			log.Println("Successfully connected to PostgreSQL")
			return &Storage{db: db}, nil
		}

		log.Printf("Ping attempt %d failed: %v", attempt, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}
