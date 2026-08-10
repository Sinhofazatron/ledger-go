package config

import (
	"fmt"
	"os"

	"perfect_api/internal/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN     string
	RedisAddr string
	SMTPHost  string
	SMTPPort  string
	SMTPFrom  string
}

func LoadConfig() *Config {
	// load file .env if there is any
	if err := godotenv.Load(); err != nil {
		logger.Log.Info("Warning: .env file not found, using environment variables")
	}

	// По умолчанию отключаем SSL (для локальной разработки)
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Формируем DSN для PostgreSQL
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "1025"
	}
	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "no-reply@gowallet.com"
	}

	return &Config{
		DBDSN:     dsn,
		RedisAddr: redisAddr,
		SMTPHost:  smtpHost,
		SMTPPort:  smtpPort,
		SMTPFrom:  smtpFrom,
	}
}
