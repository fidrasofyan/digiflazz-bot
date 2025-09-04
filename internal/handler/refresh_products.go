package handler

import (
	"context"
	"log"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/internal/job"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
)

func RefreshProducts(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		err := job.PopulateProducts(ctxWithTimeout)
		if err != nil {
			log.Printf("Error refreshing products: %v", err)

			err = service.TelegramSendMessage(ctxWithTimeout, &service.TelegramSendMessageParams{
				ChatId:    req.Message.Chat.Id,
				ParseMode: service.TelegramParseModeHTML,
				Text:      "Produk gagal diperbarui",
			})
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			return
		}

		err = service.TelegramSendMessage(ctxWithTimeout, &service.TelegramSendMessageParams{
			ChatId:    req.Message.Chat.Id,
			ParseMode: service.TelegramParseModeHTML,
			Text:      "Produk berhasil diperbarui",
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}()

	return nil, nil
}
