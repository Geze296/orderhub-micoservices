package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	HTTPPort string
	PostgresURL string
	RedisAddr string
	RedisDB int
	JWTSecret string 
}

func LoadConfig() *Config{
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Println(err)
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV","development"),
		HTTPPort: getEnv("HTTP_PORT",":8080"),
		PostgresURL: getEnv("POSTGRES_URL","postgres://postgres:postgres@localhost:5432/orderhub?sslmode=disable"),
		RedisAddr: getEnv("REDIS_ADDR","localhost:6379"),
		RedisDB: getEnvInt("REDIS_DB", 0),
		JWTSecret: getEnv("JWT_SECRET","secret"),
	}
	return cfg
}

func getEnv(key, fallback string) string{
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