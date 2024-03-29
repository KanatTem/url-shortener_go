package main

import (
	"fmt"
	"log"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	storage, err := postgres.New(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	if err := storage.Migrate(cfg.MigrationsPath); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

}
