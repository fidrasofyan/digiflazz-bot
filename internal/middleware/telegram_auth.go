package middleware

import (
	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/gofiber/fiber/v2"
)

func TelegramAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Telegram-Bot-Api-Secret-Token")
		if token != config.Cfg.TelegramWebhookSecretToken {
			return c.Status(401).SendString("Unauthorized")
		}
		return c.Next()
	}
}
