package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
)

type telegramParseMode string

const (
	TelegramParseModeHTML       telegramParseMode = "HTML"
	TelegramParseModeMarkdownV2 telegramParseMode = "MarkdownV2"
)

type TelegramSendMessageParams struct {
	ChatId             int64                             `json:"chat_id"`
	ParseMode          telegramParseMode                 `json:"parse_mode"`
	Text               string                            `json:"text"`
	LinkPreviewOptions *types.TelegramLinkPreviewOptions `json:"link_preview_options,omitempty"`
}

var telegramHttpClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

func TelegramSendMessage(ctx context.Context, params *TelegramSendMessageParams) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Cfg.TelegramBotToken)
	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := telegramHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

type TelegramAnswerCallbackQueryParams struct {
	CallbackQueryId string  `json:"callback_query_id"`
	Text            *string `json:"text,omitempty"`
	ShowAlert       *bool   `json:"show_alert,omitempty"`
}

func TelegramAnswerCallbackQuery(ctx context.Context, params *TelegramAnswerCallbackQueryParams) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", config.Cfg.TelegramBotToken)
	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := telegramHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func TelegramSetWebhook(ctx context.Context) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", config.Cfg.TelegramBotToken)
	data := struct {
		Url                string   `json:"url"`
		SecretToken        string   `json:"secret_token"`
		MaxConnections     int      `json:"max_connections"`
		DropPendingUpdates bool     `json:"drop_pending_updates"`
		AllowedUpdates     []string `json:"allowed_updates"`
	}{
		Url:                config.Cfg.WebhookURL + "/telegram",
		SecretToken:        config.Cfg.TelegramWebhookSecretToken,
		MaxConnections:     50,
		DropPendingUpdates: true,
		AllowedUpdates:     []string{"message", "callback_query"},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := telegramHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

type TelegramCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

func TelegramSetMyCommands(ctx context.Context, commands []TelegramCommand) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setMyCommands", config.Cfg.TelegramBotToken)
	data := struct {
		Commands []TelegramCommand `json:"commands"`
	}{
		Commands: commands,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := telegramHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
