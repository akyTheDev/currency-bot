package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken        string
	DatabaseURL          string
	NotificationInterval int
}

const notificationInterval int = 1

func Load() (*Config, error) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("TELEGRAM_TOKEN env variable required.")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL env variable required.")
	}

	return &Config{
		TelegramToken:        telegramToken,
		DatabaseURL:          databaseURL,
		NotificationInterval: notificationInterval,
	}, nil
}
