package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	BaseURL          string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	ShortCodeLength  int
	DefaultExpiryDays int
	RateLimitRPM     int
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		BaseURL:          getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://shortly:password@localhost:5432/shortly?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret"),
		ShortCodeLength:  getEnvInt("SHORT_CODE_LENGTH", 7),
		DefaultExpiryDays: getEnvInt("DEFAULT_EXPIRY_DAYS", 30),
		RateLimitRPM:     getEnvInt("RATE_LIMIT_RPM", 60),
	}
}

func (c *Config) DefaultExpiry() time.Duration {
	return time.Duration(c.DefaultExpiryDays) * 24 * time.Hour
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
