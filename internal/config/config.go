package config

import (
	"fmt"
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
		fmt.Println("missing env variable: APP_TIMEZONE")
		os.Exit(1)
	}
	os.Setenv("TZ", os.Getenv("APP_TIMEZONE"))

	// Telegram allowed ids
	telegramAllowedIdsStr := strings.Split(os.Getenv("TELEGRAM_ALLOWED_IDS"), ",")
	telegramAllowedIds := make([]int64, len(telegramAllowedIdsStr))
	for i := range telegramAllowedIdsStr {
		num, err := strconv.ParseInt(strings.TrimSpace(telegramAllowedIdsStr[i]), 10, 64)
		if err != nil {
			fmt.Printf("invalid TELEGRAM_ALLOWED_IDS: %s", os.Getenv("TELEGRAM_ALLOWED_IDS"))
			os.Exit(1)
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
		fmt.Println("missing env variable: APP_ENV")
		os.Exit(1)
	}
	if Cfg.AppEnv != "development" && Cfg.AppEnv != "production" {
		fmt.Printf("invalid APP_ENV: %s", Cfg.AppEnv)
		os.Exit(1)
	}
	if Cfg.AppHost == "" {
		fmt.Println("missing env variable: APP_HOST")
		os.Exit(1)
	}
	if Cfg.AppPort == "" {
		fmt.Println("missing env variable: APP_PORT")
		os.Exit(1)
	}
	if Cfg.AppName == "" {
		fmt.Println("missing env variable: APP_NAME")
		os.Exit(1)
	}
	if Cfg.TelegramBotToken == "" {
		fmt.Println("missing env variable: TELEGRAM_BOT_TOKEN")
		os.Exit(1)
	}
	if len(Cfg.TelegramAllowedIds) == 0 {
		fmt.Println("missing env variable: TELEGRAM_ALLOWED_IDS")
		os.Exit(1)
	}
	if Cfg.DigiflazzBaseUrl == "" {
		fmt.Println("missing env variable: DIGIFLAZZ_BASE_URL")
		os.Exit(1)
	}
	if Cfg.DigiflazzUsername == "" {
		fmt.Println("missing env variable: DIGIFLAZZ_USERNAME")
		os.Exit(1)
	}
	if Cfg.DigiflazzApiKey == "" {
		fmt.Println("missing env variable: DIGIFLAZZ_API_KEY")
		os.Exit(1)
	}
	if Cfg.DatabaseURL == "" {
		fmt.Println("missing env variable: DATABASE_URL")
		os.Exit(1)
	}
	if Cfg.WebhookURL == "" {
		fmt.Println("missing env variable: WEBHOOK_URL")
		os.Exit(1)
	}
	if Cfg.TelegramWebhookSecretToken == "" {
		fmt.Println("missing env variable: TELEGRAM_WEBHOOK_SECRET_TOKEN")
		os.Exit(1)
	}
	if Cfg.DigiflazzWebhookSecretToken == "" {
		fmt.Println("missing env variable: DIGIFLAZZ_WEBHOOK_SECRET_TOKEN")
		os.Exit(1)
	}

	// Generate Telegram webhook secret
	if Cfg.TelegramWebhookSecretToken == "" || Cfg.TelegramWebhookSecretToken == "auto" {
		secretToken, err := util.GenerateSecretToken(32)
		if err != nil {
			fmt.Printf("failed to generate secret token: %v", err)
			os.Exit(1)
		}
		Cfg.TelegramWebhookSecretToken = secretToken
	}
}
