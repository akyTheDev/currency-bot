package config

import (
	"os"
	"testing"
)

const dummyDBURL string = "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
const dummyTelegramToken string = "dummy_token"

func TestLoad_Success(t *testing.T) {
	t.Setenv("TELEGRAM_TOKEN", dummyTelegramToken)
	t.Setenv("DATABASE_URL", dummyDBURL)

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Unexpected error :%v", err)
	}

	if cfg.TelegramToken != dummyTelegramToken {
		t.Errorf("TelegramToken=%q, expected %q", cfg.TelegramToken, dummyTelegramToken)
	}

	if cfg.DatabaseURL != dummyDBURL {
		t.Errorf("DatabaseURL=%q, expected %q", cfg.DatabaseURL, dummyDBURL)
	}

	if cfg.NotificationInterval != notificationInterval {
		t.Errorf("NotificationInterval=%q, expected %q", cfg.NotificationInterval, notificationInterval)
	}
}

func TestLoad_MissingEnv(t *testing.T) {
	os.Unsetenv("TELEGRAM_TOKEN")
	os.Unsetenv("DATABASE_URL")

	tests := []struct {
		name            string
		setEnv          func()
		missingVarError string
	}{
		{
			name: "Missing TELEGRAM_TOKEN",
			setEnv: func() {
				os.Unsetenv("TELEGRAM_TOKEN")
				os.Setenv("DATABASE_URL", dummyDBURL)
			},
			missingVarError: "TELEGRAM_TOKEN env variable required.",
		},
		{
			name: "Missing DATABASE_URL",
			setEnv: func() {
				os.Unsetenv("DATABASE_URL")
				os.Setenv("TELEGRAM_TOKEN", dummyTelegramToken)
			},
			missingVarError: "DATABASE_URL env variable required.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setEnv()

			cfg, err := Load()
			if err == nil {
				t.Fatalf("Error wasn't returned. cfg: %+v", cfg)
			}

			if err.Error() != tc.missingVarError {
				t.Errorf("Error: %v, Expected Error: %v", err.Error(), tc.missingVarError)
			}
		})
	}
}
