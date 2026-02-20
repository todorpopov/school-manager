package configs

import (
	"os"
)

type Config struct {
	ApiPort string
}

func ParseConfig() *Config {
	return &Config{
		ApiPort: GetEnv("API_PORT", "8080"),
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
