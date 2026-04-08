package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	DBDSN       string
	AppURL      string
}

func Load() *Config {
	env := os.Getenv("ENVIRONMENT")

	// Load .env ONLY for local/dev
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (ok in production)")
		}
	}

	cfg := &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		DBDSN:       os.Getenv("DB_DSN"),
		AppURL:      os.Getenv("APP_URL"),
	}

	return cfg
}