package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DB_DSN string
}

func LoadConfig() (*Config, error) {
    _ = godotenv.Load() // Load từ file .env, nếu không có thì đọc biến env hệ thống

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSLMODE"),
    )

    return &Config{
        DB_DSN: dsn,
    }, nil
}
