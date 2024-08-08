package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host      string
	Port      string
	DBURL     string
	MarketURL string
	//JWTExpirationInSeconds int64
	//JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:      getEnv("HOST", "http://localhost"),
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", ""),
		MarketURL: getEnv("MARKET_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
