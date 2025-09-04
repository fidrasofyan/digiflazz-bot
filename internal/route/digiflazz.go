package route

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
	"github.com/gofiber/fiber/v2"
)

func Digiflazz() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req types.DigiflazzUpdate
		if err := c.BodyParser(&req); err != nil {
			return util.NewError(err)
		}

		go func() {
			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			var textB strings.Builder
			textB.WriteString(fmt.Sprintf(
				"%s ke %s %s. SN: <code>%s</code>. ",
				req.Data.BuyerSKUCode,
				req.Data.CustomerNo,
				req.Data.Status,
				*req.Data.SN,
			))
			textB.WriteString(util.Sprintf("Harga: %d. Saldo: %d. ", req.Data.Price, req.Data.BuyerLastSaldo))
			textB.WriteString(fmt.Sprintf("Waktu: %s. ", time.Now().Format("2 Jan 2006 15:04:05 MST")))
			textB.WriteString(fmt.Sprintf("Keterangan: %s", req.Data.Message))

			for _, chatId := range config.Cfg.TelegramAllowedIds {
				err := service.TelegramSendMessage(ctxWithTimeout, &service.TelegramSendMessageParams{
					ChatId:    chatId,
					ParseMode: service.TelegramParseModeHTML,
					Text:      textB.String(),
				})
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
		}()

		return c.Status(200).SendString("OK")
	}
}
