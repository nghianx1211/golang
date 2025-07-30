package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Load from .env if has

	cfg := &Config{}

	cfg.Server = ServerConfig{
		Port: getEnv("SERVER_PORT", "8080"),
	}

	cfg.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "user_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	cfg.JWT = JWTConfig{
		Secret: getEnv("JWT_SECRET", "super-secret-key"),
		RefreshSecret: getEnv("JWT_REFRESH_SECRET", "super-secret-key"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}