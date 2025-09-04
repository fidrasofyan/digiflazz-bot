package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/database"
	"github.com/fidrasofyan/digiflazz-bot/database/repository"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
	"github.com/google/uuid"
)

var trxCmd = "_transaction"

type trxData struct {
	Code   string `json:"code"`
	Number string `json:"number"`
}

func Transaction(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	var chatId int64

	// Is it callback query?
	if req.CallbackQuery != nil {
		chatId = req.CallbackQuery.From.Id
	} else {
		chatId = req.Message.Chat.Id
	}

	// Get chat
	chat, err := repository.TelegramGetChat(ctx, chatId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, util.NewError(err)
	}

	if chat.ID == 0 {
		// Create new chat
		chat, err = repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
			ID:      chatId,
			Command: trxCmd,
			Step:    1,
		})
		if err != nil {
			return nil, util.NewError(err)
		}
	}

	switch chat.Step {
	// Step 1
	case 1:
		textParts := strings.Fields(req.Message.Text)
		productCode := textParts[0]
		destinationNumber := textParts[1]

		// Prepaid product exist?
		prepaidProduct, err := database.Sqlc.GetPrepaidProductBySKUCode(ctx, productCode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewError(err)
		}
		if prepaidProduct.ID == 0 {
			// Delete step
			err := repository.TelegramDeleteChat(ctx, chatId)
			if err != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      req.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "<i>Produk tidak ditemukan</i>",
				ReplyMarkup: types.DefaultReplyMarkup,
			}, nil
		}

		var prepaidProductStatus string
		if prepaidProduct.BuyerProductStatus && prepaidProduct.SellerProductStatus {
			prepaidProductStatus = "✅"
		} else {
			prepaidProductStatus = "❌"
		}

		var textB strings.Builder
		textB.WriteString(fmt.Sprintf("Kode: %s\n", prepaidProduct.BuyerSkuCode))
		textB.WriteString(fmt.Sprintf("Tujuan: %s\n", destinationNumber))
		textB.WriteString(util.Sprintf("Harga: Rp %d\n\n", prepaidProduct.Price))
		textB.WriteString(fmt.Sprintf("Seller: %s\n", prepaidProduct.SellerName))
		textB.WriteString(fmt.Sprintf("Status: %s\n", prepaidProductStatus))
		textB.WriteString(fmt.Sprintf("Nama: %s\n", prepaidProduct.Name))
		textB.WriteString(fmt.Sprintf("Deskripsi: %s\n", *prepaidProduct.Description))
		textB.WriteString("\nYakin ingin memproses?")

		// Set step
		trxData := &trxData{
			Code:   productCode,
			Number: destinationNumber,
		}
		trxDataB, err := json.Marshal(trxData)
		if err != nil {
			return nil, util.NewError(err)
		}
		_, err = repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
			ID:      chatId,
			Command: trxCmd,
			Step:    2,
			Data:    trxDataB,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodSendMessage,
			ChatId:    req.Message.Chat.Id,
			ParseMode: types.TelegramParseModeHTML,
			Text:      textB.String(),
			ReplyMarkup: types.TelegramReplyKeyboardMarkup{
				ResizeKeyboard: true,
				Keyboard: [][]string{
					{"Ya", "Tidak"},
				},
			},
		}, nil

	// Step 2
	case 2:
		// Delete step
		defer func() {
			err = repository.TelegramDeleteChat(ctx, chatId)
			if err != nil {
				log.Println(err)
			}
		}()

		if req.Message.Text != "Ya" {
			return &types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      req.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "<i>Dibatalkan</i>",
				ReplyMarkup: types.DefaultReplyMarkup,
			}, nil
		}

		trxData := &trxData{}
		err := json.Unmarshal(chat.Data, trxData)
		if err != nil {
			return nil, util.NewError(err)
		}

		// Send to digiflazz
		digiflazzRes, err := service.DigiflazzCreateTrx(ctx, &service.DigiflazzCreateTrxParams{
			RefID:        uuid.Must(uuid.NewV7()).String(),
			BuyerSKUCode: trxData.Code,
			CustomerNo:   trxData.Number,
		})
		if err != nil {
			var digiflazzError *service.DigiflazzErrorResponse
			if errors.As(err, &digiflazzError) {
				return &types.TelegramResponse{
					Method:      types.TelegramMethodSendMessage,
					ChatId:      req.Message.Chat.Id,
					ParseMode:   types.TelegramParseModeHTML,
					Text:        digiflazzError.Error(),
					ReplyMarkup: types.DefaultReplyMarkup,
				}, nil
			}
			return nil, util.NewError(err)
		}

		if digiflazzRes.Data.RC == "03" || digiflazzRes.Data.RC == "99" {
			return &types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      req.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        fmt.Sprintf("<i>%s ke %s sedang diproses...</i>", trxData.Code, trxData.Number),
				ReplyMarkup: types.DefaultReplyMarkup,
			}, nil
		}

		var textB strings.Builder
		textB.WriteString(fmt.Sprintf(
			"%s ke %s %s. SN: <code>%s</code>. ",
			digiflazzRes.Data.BuyerSKUCode,
			digiflazzRes.Data.CustomerNo,
			digiflazzRes.Data.Status,
			*digiflazzRes.Data.SN,
		))
		textB.WriteString(fmt.Sprintf("Waktu: %s. ", time.Now().Format("2 Jan 2006 15:04:05 MST")))
		textB.WriteString(fmt.Sprintf("Keterangan: %s", digiflazzRes.Data.Message))

		return &types.TelegramResponse{
			Method:      types.TelegramMethodSendMessage,
			ChatId:      req.Message.Chat.Id,
			ParseMode:   types.TelegramParseModeHTML,
			Text:        textB.String(),
			ReplyMarkup: types.DefaultReplyMarkup,
		}, nil

	// Unhandled step
	default:
		// Delete step
		err := repository.TelegramDeleteChat(ctx, chatId)
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:      types.TelegramMethodSendMessage,
			ChatId:      req.Message.Chat.Id,
			ParseMode:   types.TelegramParseModeHTML,
			Text:        "<i>Unhandled step</i>",
			ReplyMarkup: types.DefaultReplyMarkup,
		}, nil
	}

}
