// internal/config/config.go
package config // Add package declaration

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string           `yaml:"env" env-default:"development"`
	Postgres       PostgreSQLConfig `yaml:"postgres" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host" env:"PGHOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PGPORT" env-default:"5432"`
	User     string `yaml:"user" env:"PGUSER" env-default:"postgres"`
	Password string `yaml:"password" env:"PGPASSWORD"`
	DBName   string `yaml:"dbname" env:"PGDATABASE"`
	SSLMode  string `yaml:"sslmode" env:"PGSSLMODE" env-default:"disable"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
