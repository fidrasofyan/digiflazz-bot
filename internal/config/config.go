package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fidrasofyan/digiflazz-bot/internal/util"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                      string
	AppHost                     string
	AppPort                     string
	AppName                     string
	TelegramBotToken            string
	TelegramAllowedIds          []int64
	DigiflazzBaseUrl            string
	DigiflazzUsername           string
	DigiflazzApiKey             string
	DatabaseURL                 string
	WebhookURL                  string
	TelegramWebhookSecretToken  string
	DigiflazzWebhookSecretToken string
}

var Cfg *Config

func MustLoadConfig() {
	godotenv.Load()

	// App timezone
	if os.Getenv("APP_TIMEZONE") == "" {
		log.Fatal("missing APP_TIMEZONE")
	}
	os.Setenv("TZ", os.Getenv("APP_TIMEZONE"))

	// Telegram allowed ids
	telegramAllowedIdsStr := strings.Split(os.Getenv("TELEGRAM_ALLOWED_IDS"), ",")
	telegramAllowedIds := make([]int64, len(telegramAllowedIdsStr))
	for i := range telegramAllowedIdsStr {
		num, err := strconv.ParseInt(strings.TrimSpace(telegramAllowedIdsStr[i]), 10, 64)
		if err != nil {
			log.Fatalf("invalid TELEGRAM_ALLOWED_IDS: %s", os.Getenv("TELEGRAM_ALLOWED_IDS"))
		}
		telegramAllowedIds[i] = num
	}

	Cfg = &Config{
		AppEnv:                      os.Getenv("APP_ENV"),
		AppHost:                     os.Getenv("APP_HOST"),
		AppPort:                     os.Getenv("APP_PORT"),
		AppName:                     os.Getenv("APP_NAME"),
		TelegramBotToken:            os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramAllowedIds:          telegramAllowedIds,
		DigiflazzBaseUrl:            os.Getenv("DIGIFLAZZ_BASE_URL"),
		DigiflazzUsername:           os.Getenv("DIGIFLAZZ_USERNAME"),
		DigiflazzApiKey:             os.Getenv("DIGIFLAZZ_API_KEY"),
		DatabaseURL:                 os.Getenv("DATABASE_URL"),
		WebhookURL:                  os.Getenv("WEBHOOK_URL"),
		TelegramWebhookSecretToken:  os.Getenv("TELEGRAM_WEBHOOK_SECRET_TOKEN"),
		DigiflazzWebhookSecretToken: os.Getenv("DIGIFLAZZ_WEBHOOK_SECRET_TOKEN"),
	}

	// Validate
	if Cfg.AppEnv == "" {
		log.Fatal("missing APP_ENV")
	}
	if Cfg.AppEnv != "development" && Cfg.AppEnv != "production" {
		log.Fatalf("invalid APP_ENV: %s", Cfg.AppEnv)
	}
	if Cfg.AppHost == "" {
		log.Fatal("missing APP_HOST")
	}
	if Cfg.AppPort == "" {
		log.Fatal("missing APP_PORT")
	}
	if Cfg.AppName == "" {
		log.Fatal("missing APP_NAME")
	}
	if Cfg.TelegramBotToken == "" {
		log.Fatal("missing TELEGRAM_BOT_TOKEN")
	}
	if len(Cfg.TelegramAllowedIds) == 0 {
		log.Fatal("missing TELEGRAM_ALLOWED_IDS")
	}
	if Cfg.DigiflazzBaseUrl == "" {
		log.Fatal("missing DIGIFLAZZ_BASE_URL")
	}
	if Cfg.DigiflazzUsername == "" {
		log.Fatal("missing DIGIFLAZZ_USERNAME")
	}
	if Cfg.DigiflazzApiKey == "" {
		log.Fatal("missing DIGIFLAZZ_API_KEY")
	}
	if Cfg.DatabaseURL == "" {
		log.Fatal("missing DATABASE_URL")
	}
	if Cfg.WebhookURL == "" {
		log.Fatal("missing WEBHOOK_URL")
	}
	if Cfg.TelegramWebhookSecretToken == "" {
		log.Fatal("missing TELEGRAM_WEBHOOK_SECRET_TOKEN")
	}
	if Cfg.DigiflazzWebhookSecretToken == "" {
		log.Fatal("missing DIGIFLAZZ_WEBHOOK_SECRET_TOKEN")
	}

	// Generate Telegram webhook secret
	if Cfg.TelegramWebhookSecretToken == "" || Cfg.TelegramWebhookSecretToken == "auto" {
		secretToken, err := util.GenerateSecretToken(32)
		if err != nil {
			log.Fatalf("failed to generate secret token: %v", err)
		}
		Cfg.TelegramWebhookSecretToken = secretToken
	}
}
