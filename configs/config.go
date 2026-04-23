package configs

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ApiPort         string
	SystemAuthToken string

	DBUrl                   string
	DBMinConnections        int32
	DBMaxConnections        int32
	DBMaxConnectionLifetime time.Duration
	DBMaxConnectionIdleTime time.Duration

	SessionExpiration time.Duration
}

func ParseConfig() *Config {
	return &Config{
		ApiPort:                 GetEnv("API_PORT", "8080"),
		SystemAuthToken:         GetEnv("SYSTEM_AUTH_TOKEN", "admin"),
		DBUrl:                   GetEnv("DB_URL", ""),
		DBMinConnections:        GetEnvInt32("DB_MIN_CONNECTIONS", 5),
		DBMaxConnections:        GetEnvInt32("DB_MAX_CONNECTIONS", 10),
		DBMaxConnectionLifetime: GetEnvDuration("DB_MAX_CONNECTION_LIFETIME", 5*time.Minute),
		DBMaxConnectionIdleTime: GetEnvDuration("DB_MAX_CONNECTION_IDLE_TIME", 10*time.Minute),
		SessionExpiration:       GetEnvDuration("SESSION_EXPIRATION", time.Hour),
	}
}

func NewTestConfig() *Config {
	return &Config{
		ApiPort:                 "8080",
		DBUrl:                   "postgres://admin:admin@100.109.199.26:30001/school_manager",
		DBMinConnections:        5,
		DBMaxConnections:        10,
		DBMaxConnectionLifetime: 5 * time.Minute,
		DBMaxConnectionIdleTime: 10 * time.Minute,
		SessionExpiration:       time.Hour,
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt32(key string, fallback int32) int32 {
	if value, ok := os.LookupEnv(key); ok {
		ret, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fallback
		}

		return int32(ret)
	}

	return fallback
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		ret, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		return ret
	}
	return fallback
}
