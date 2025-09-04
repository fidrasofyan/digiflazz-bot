package route

import (
	"database/sql"
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/fidrasofyan/digiflazz-bot/database/repository"
	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/handler"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
	"github.com/gofiber/fiber/v2"
)

var trxRegex = regexp.MustCompile(`^([A-Za-z0-9-]+)\s+(\d+)$`)

func Telegram() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.TelegramUpdate
		if err := c.BodyParser(&req); err != nil {
			return util.NewError(err)
		}

		var chatId int64
		var command string

		// Is it callback query?
		if req.CallbackQuery != nil {
			// Set chat id
			chatId = req.CallbackQuery.From.Id
		} else {
			// Set chat id
			chatId = req.Message.Chat.Id

			// Only text message is supported
			if req.Message.Text == "" {
				return c.Status(200).JSON(types.TelegramResponse{
					Method:      types.TelegramMethodSendMessage,
					ChatId:      req.Message.Chat.Id,
					ParseMode:   types.TelegramParseModeHTML,
					Text:        "<i>Only text command is supported</i>",
					ReplyMarkup: types.DefaultReplyMarkup,
				})
			}

			// Set command
			command = strings.TrimSpace(strings.ToLower(
				req.Message.Text,
			))
			// Limit command length
			if len(command) > 30 {
				command = command[:30]
			}
			// Remove leading slashes
			command = strings.TrimLeft(command, "/")
			// Remove leading underscores to prevent calling internal commands
			command = strings.TrimLeft(command, "_")
		}

		// Is it "cancel" command?
		if command == "cancel" {
			// Delete chat
			err := repository.TelegramDeleteChat(c.UserContext(), chatId)
			if err != nil {
				return util.NewError(err)
			}
			return c.Status(200).JSON(types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      req.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "<i>Dibatalkan</i>",
				ReplyMarkup: types.DefaultReplyMarkup,
			})
		}

		// Get chat
		chat, err := repository.TelegramGetChat(c.UserContext(), chatId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return util.NewError(err)
		}
		if chat.ID != 0 {
			// Set command
			command = chat.Command
		}

		// Is chat ID allowed?
		if command != "start" && !slices.Contains(config.Cfg.TelegramAllowedIds, chatId) {
			return c.Status(200).JSON(types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      req.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "<i>Access denied</i>",
				ReplyMarkup: types.DefaultReplyMarkup,
			})
		}

		switch command {
		// Start
		case "start":
			resp, err := handler.Start(c.UserContext(), &req)
			if err != nil {
				return util.NewError(err)
			}
			if resp == nil {
				return c.Status(200).SendString("OK")
			}
			return c.Status(200).JSON(resp)

		// Product list
		case "daftar produk":
			resp, err := handler.ProductList(c.UserContext(), &req)
			if err != nil {
				return util.NewError(err)
			}
			if resp == nil {
				return c.Status(200).SendString("OK")
			}
			return c.Status(200).JSON(resp)

		// Refresh products
		case "refresh produk":
			resp, err := handler.RefreshProducts(c.UserContext(), &req)
			if err != nil {
				return util.NewError(err)
			}
			if resp == nil {
				return c.Status(200).SendString("OK")
			}
			return c.Status(200).JSON(resp)

		// Check balance
		case "cek saldo":
			resp, err := handler.CheckBalance(c.UserContext(), &req)
			if err != nil {
				return util.NewError(err)
			}
			if resp == nil {
				return c.Status(200).SendString("OK")
			}
			return c.Status(200).JSON(resp)

		// Transaction
		case "_transaction":
			resp, err := handler.Transaction(c.UserContext(), &req)
			if err != nil {
				return util.NewError(err)
			}
			if resp == nil {
				return c.Status(200).SendString("OK")
			}
			return c.Status(200).JSON(resp)

		// Not found
		default:
			// Is it transaction?
			if req.Message != nil {
				req.Message.Text = strings.TrimSpace(req.Message.Text)
				req.Message.Text = strings.ReplaceAll(req.Message.Text, "-", "")
				req.Message.Text = strings.ReplaceAll(req.Message.Text, "+62", "0")
				if trxRegex.MatchString(req.Message.Text) {
					resp, err := handler.Transaction(c.UserContext(), &req)
					if err != nil {
						return util.NewError(err)
					}
					if resp == nil {
						return c.Status(200).SendString("OK")
					}
					return c.Status(200).JSON(resp)
				}
			}

			resp, err := handler.NotFound(c.UserContext(), &req)
			if err != nil || resp == nil {
				return util.NewError(err)
			}
			return c.Status(200).JSON(resp)
		}
	}
}
