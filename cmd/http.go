package cmd

import (
	"errors"
	"log"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/database/repository"
	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/middleware"
	"github.com/fidrasofyan/digiflazz-bot/internal/route"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func MustStartHTTPServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Digiflazz Bot",
		Prefork:               false,
		DisableStartupMessage: true,
		BodyLimit:             10 * 1024, // 10KB
		ReadBufferSize:        4 * 1024,  // 4KB
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
			}

			var body types.TelegramUpdate
			if err := c.BodyParser(&body); err != nil {
				err = c.Status(code).JSON(&fiber.Map{
					"ok":      false,
					"message": fiberErr.Message,
				})

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"ok":      false,
						"message": "Internal server error	",
					})
				}

				return nil
			}

			// At this point, the body is a valid Telegram update
			log.Printf("Error: %v", err)

			text := "<i>Something went wrong</i>"
			if errors.Is(err, fiber.ErrRequestTimeout) {
				text = "<i>Request timeout</i>"
			}

			// Delete chat
			_ = repository.TelegramDeleteChat(c.Context(), body.CallbackQuery.From.Id)

			if body.CallbackQuery != nil {
				// Answer callback query
				_ = service.TelegramAnswerCallbackQuery(c.Context(), &service.TelegramAnswerCallbackQueryParams{
					CallbackQueryId: body.CallbackQuery.Id,
				})

				return c.Status(200).JSON(types.TelegramResponse{
					Method:    types.TelegramMethodEditMessageText,
					MessageId: body.CallbackQuery.Message.MessageId,
					ChatId:    body.CallbackQuery.Message.Chat.Id,
					ParseMode: types.TelegramParseModeHTML,
					Text:      text,
				})
			}

			return c.Status(200).JSON(types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      body.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        text,
				ReplyMarkup: types.DefaultReplyMarkup,
			})
		},
	})

	// Middlewares
	app.Use(recover.New())
	if config.Cfg.AppEnv == "development" {
		app.Use(logger.New())
	}

	// Routes
	app.Post(
		"/telegram",
		middleware.TelegramAuth(),
		timeout.NewWithContext(route.Telegram(), 10*time.Second),
	)
	app.Post(
		"/digiflazz",
		middleware.DigiflazzAuth(),
		timeout.NewWithContext(route.Digiflazz(), 10*time.Second),
	)

	// Not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("404 not found")
	})

	go func() {
		log.Printf("Server is running on http://%s:%s", config.Cfg.AppHost, config.Cfg.AppPort)
		err := app.Listen(config.Cfg.AppHost + ":" + config.Cfg.AppPort)
		if err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	return app
}
