package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	MAX_BUFFER_SIZE   = 5_000
	MAX_CHUNK_SIZE    = 100
	MAX_FETCH_DB      = 100_000
	MAX_FETCH_STREAMS = 15_000
	MAX_PER_WORKER    = 2_500_000
)

type Config struct {
	PgHost     string
	PgPort     string
	PgUser     string
	PgPassword string
	PgDatabase string
	PgSchema   string
	ServerPort string
}

func init() {
	godotenv.Load()
}

func Get() *Config {
	return &Config{
		PgHost:     os.Getenv("DB_HOST"),
		PgPort:     os.Getenv("DB_PORT"),
		PgUser:     os.Getenv("DB_USER"),
		PgPassword: os.Getenv("DB_PASSWORD"),
		PgDatabase: os.Getenv("DB_DATABASE"),
		PgSchema:   os.Getenv("DB_SCHEMA"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}
