package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	PostgresURL string
	RedisAddr   string
	RedisDB     int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/orderhub?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:     getEnvInt("REDIS_DB", 0),
	}
	
	return cfg, nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	if val == 0 {
		return fallback
	}
	return val
}
