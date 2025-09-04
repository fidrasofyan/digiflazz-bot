package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/gofiber/fiber/v2"
)

func DigiflazzAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		signature := c.Get("X-Hub-Signature")
		if signature == "" {
			return c.Status(401).SendString("Unauthorized")
		}

		event := c.Get("X-Digiflazz-Event")
		if event != "update" {
			return c.Status(200).SendString("OK")
		}

		// Verify signature
		mac := hmac.New(sha1.New, []byte(config.Cfg.DigiflazzWebhookSecretToken))
		mac.Write(c.Body())
		expectedMAC := mac.Sum(nil)
		expectedSignature := "sha1=" + hex.EncodeToString(expectedMAC)

		valid := hmac.Equal([]byte(expectedSignature), []byte(signature))
		if !valid {
			return c.Status(401).SendString("Unauthorized")
		}

		return c.Next()
	}
}
