package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/database"
	"github.com/fidrasofyan/digiflazz-bot/database/repository"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
)

var productListCmd = "daftar produk"

func ProductList(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
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
			Command: productListCmd,
			Step:    1,
		})
		if err != nil {
			return nil, util.NewError(err)
		}
	}

	switch chat.Step {
	// Step 1
	case 1:
		// Get categories
		categories, err := database.Sqlc.GetCategories(ctx)
		if err != nil {
			return nil, util.NewError(err)
		}

		if len(categories) == 0 {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:      types.TelegramMethodSendMessage,
				ChatId:      chatId,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "<i>Tidak ada produk</i>",
				ReplyMarkup: types.DefaultReplyMarkup,
			}, nil
		}

		inlineKeyboard := make([][]types.TelegramInlineKeyboardButton, len(categories)+1)

		for i, category := range categories {
			inlineKeyboard[i] = []types.TelegramInlineKeyboardButton{
				{
					Text:         category,
					CallbackData: category,
				},
			}
		}

		inlineKeyboard[len(categories)] = []types.TelegramInlineKeyboardButton{
			{
				Text: "❌", CallbackData: "cancel",
			},
		}

		// Set step
		_, err = repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
			ID:      chatId,
			Command: productListCmd,
			Step:    2,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodSendMessage,
			ChatId:    chatId,
			ParseMode: types.TelegramParseModeHTML,
			Text:      "Pilih kategori:",
			ReplyMarkup: types.TelegramInlineKeyboardMarkup{
				InlineKeyboard: inlineKeyboard,
			},
		}, nil

	// Step 2
	case 2:
		// It must be callback query
		if req.CallbackQuery == nil {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodSendMessage,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Perintah tidak valid</i>",
			}, nil
		}

		// Answer callback query
		go func() {
			acqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			service.TelegramAnswerCallbackQuery(acqCtx, &service.TelegramAnswerCallbackQueryParams{
				CallbackQueryId: req.CallbackQuery.Id,
			})
		}()

		// Cancel
		if req.CallbackQuery.Data == "cancel" {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodEditMessageText,
				MessageId: req.CallbackQuery.Message.MessageId,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Dibatalkan</i>",
			}, nil
		}

		// Get brands
		category := req.CallbackQuery.Data
		brands, err := database.Sqlc.GetBrandsByCategory(ctx, category)
		if err != nil {
			return nil, util.NewError(err)
		}

		inlineKeyboard := make([][]types.TelegramInlineKeyboardButton, len(brands)+1)

		for i, brand := range brands {
			inlineKeyboard[i] = []types.TelegramInlineKeyboardButton{
				{
					Text:         brand,
					CallbackData: fmt.Sprintf("%s,%s", category, brand),
				},
			}
		}

		inlineKeyboard[len(brands)] = []types.TelegramInlineKeyboardButton{
			{
				Text: "⬅️", CallbackData: "back",
			},
			{
				Text: "❌", CallbackData: "cancel",
			},
		}

		// Set previous reply markup
		previousReplyMarkupB, err := json.Marshal(req.CallbackQuery.Message.ReplyMarkup)
		if err != nil {
			return nil, util.NewError(err)
		}
		err = repository.TelegramSetReplyMarkup(ctx, &repository.TelegramSetReplyMarkupParams{
			ID:          chatId,
			Step:        1,
			ReplyMarkup: previousReplyMarkupB,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		// Set step
		_, err = repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
			ID:      chatId,
			Command: productListCmd,
			Step:    3,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodEditMessageText,
			MessageId: req.CallbackQuery.Message.MessageId,
			ChatId:    req.CallbackQuery.Message.Chat.Id,
			ParseMode: types.TelegramParseModeHTML,
			Text:      "Pilih provider:",
			ReplyMarkup: types.TelegramInlineKeyboardMarkup{
				InlineKeyboard: inlineKeyboard,
			},
		}, nil

	// Step 3
	case 3:
		// It must be callback query
		if req.CallbackQuery == nil {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodSendMessage,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Perintah tidak valid</i>",
			}, nil
		}

		// Answer callback query
		go func() {
			acqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			service.TelegramAnswerCallbackQuery(acqCtx, &service.TelegramAnswerCallbackQueryParams{
				CallbackQueryId: req.CallbackQuery.Id,
			})
		}()

		// Cancel
		if req.CallbackQuery.Data == "cancel" {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodEditMessageText,
				MessageId: req.CallbackQuery.Message.MessageId,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Dibatalkan</i>",
			}, nil
		}

		// Back
		if req.CallbackQuery.Data == "back" {
			// Set step
			_, err := repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
				ID:      chatId,
				Command: productListCmd,
				Step:    2,
			})
			if err != nil {
				return nil, util.NewError(err)
			}

			replyMarkup := &types.TelegramInlineKeyboardMarkup{}
			err = json.Unmarshal(chat.ReplyMarkup1, replyMarkup)
			if err != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:      types.TelegramMethodEditMessageText,
				MessageId:   req.CallbackQuery.Message.MessageId,
				ChatId:      req.CallbackQuery.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "Pilih kategori:",
				ReplyMarkup: replyMarkup,
			}, nil
		}

		data := strings.Split(req.CallbackQuery.Data, ",")
		productCategory := data[0]
		productBrand := data[1]

		// Get product types
		productTypes, err := database.Sqlc.GetTypesByCategoryAndBrand(ctx, &database.GetTypesByCategoryAndBrandParams{
			Category: productCategory,
			Brand:    productBrand,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		inlineKeyboard := make([][]types.TelegramInlineKeyboardButton, len(productTypes)+1)

		for i, pt := range productTypes {
			inlineKeyboard[i] = []types.TelegramInlineKeyboardButton{
				{
					Text:         pt,
					CallbackData: fmt.Sprintf("%s,%s,%s", productCategory, productBrand, pt),
				},
			}
		}

		inlineKeyboard[len(productTypes)] = []types.TelegramInlineKeyboardButton{
			{
				Text: "⬅️", CallbackData: "back",
			},
			{
				Text: "❌", CallbackData: "cancel",
			},
		}

		// Set previous reply markup
		previousReplyMarkupB, err := json.Marshal(req.CallbackQuery.Message.ReplyMarkup)
		if err != nil {
			return nil, util.NewError(err)
		}
		err = repository.TelegramSetReplyMarkup(ctx, &repository.TelegramSetReplyMarkupParams{
			ID:          chatId,
			Step:        2,
			ReplyMarkup: previousReplyMarkupB,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		// Set step
		_, err = repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
			ID:      chatId,
			Command: productListCmd,
			Step:    4,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodEditMessageText,
			MessageId: req.CallbackQuery.Message.MessageId,
			ChatId:    req.CallbackQuery.Message.Chat.Id,
			ParseMode: types.TelegramParseModeHTML,
			Text:      "Pilih tipe:",
			ReplyMarkup: types.TelegramInlineKeyboardMarkup{
				InlineKeyboard: inlineKeyboard,
			},
		}, nil

	// Step 4
	case 4:
		// It must be callback query
		if req.CallbackQuery == nil {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodSendMessage,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Perintah tidak valid</i>",
			}, nil
		}

		// Answer callback query
		go func() {
			acqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			service.TelegramAnswerCallbackQuery(acqCtx, &service.TelegramAnswerCallbackQueryParams{
				CallbackQueryId: req.CallbackQuery.Id,
			})
		}()

		// Cancel
		if req.CallbackQuery.Data == "cancel" {
			// Delete chat
			if repository.TelegramDeleteChat(ctx, chatId) != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:    types.TelegramMethodEditMessageText,
				MessageId: req.CallbackQuery.Message.MessageId,
				ChatId:    chatId,
				ParseMode: types.TelegramParseModeHTML,
				Text:      "<i>Dibatalkan</i>",
			}, nil
		}

		// Back
		if req.CallbackQuery.Data == "back" {
			// Set step
			_, err := repository.TelegramSetChat(ctx, &repository.TelegramSetChatParams{
				ID:      chatId,
				Command: productListCmd,
				Step:    3,
			})
			if err != nil {
				return nil, util.NewError(err)
			}

			replyMarkup := &types.TelegramInlineKeyboardMarkup{}
			err = json.Unmarshal(chat.ReplyMarkup2, replyMarkup)
			if err != nil {
				return nil, util.NewError(err)
			}

			return &types.TelegramResponse{
				Method:      types.TelegramMethodEditMessageText,
				MessageId:   req.CallbackQuery.Message.MessageId,
				ChatId:      req.CallbackQuery.Message.Chat.Id,
				ParseMode:   types.TelegramParseModeHTML,
				Text:        "Pilih provider:",
				ReplyMarkup: replyMarkup,
			}, nil
		}

		data := strings.Split(req.CallbackQuery.Data, ",")
		productCategory := data[0]
		productBrand := data[1]
		productType := data[2]

		// Get prepaid products
		prepaidProducts, err := database.Sqlc.GetPrepaidProducts(ctx, &database.GetPrepaidProductsParams{
			Category: productCategory,
			Brand:    productBrand,
			Type:     productType,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		textLimit := 3500
		var textB strings.Builder
		textB.WriteString(fmt.Sprintf("<b>%s » %s » %s</b>\n\n", productCategory, productBrand, productType))

		for _, pp := range prepaidProducts {
			var status string
			if pp.BuyerProductStatus && pp.SellerProductStatus {
				status = "✅"
			} else {
				status = "❌"
			}
			textB.WriteString(fmt.Sprintf("%s Kode: <code>%s</code>\n", status, pp.BuyerSkuCode))
			textB.WriteString(fmt.Sprintf("Nama: %s\n", pp.Name))
			textB.WriteString(fmt.Sprintf("Seller: %s\n", pp.SellerName))
			textB.WriteString(util.Sprintf("Harga: Rp %d\n", pp.Price))

			// If text is too long, send it part by part
			if textB.Len() >= textLimit {
				service.TelegramSendMessage(ctx, &service.TelegramSendMessageParams{
					ChatId:    req.CallbackQuery.Message.Chat.Id,
					ParseMode: service.TelegramParseModeHTML,
					Text:      textB.String(),
					LinkPreviewOptions: &types.TelegramLinkPreviewOptions{
						IsDisabled: true,
					},
				})
				textB.Reset()
			}
		}

		// Delete chat
		if repository.TelegramDeleteChat(ctx, chatId) != nil {
			return nil, util.NewError(err)
		}

		if textB.Len() == 0 {
			return nil, nil
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodEditMessageText,
			MessageId: req.CallbackQuery.Message.MessageId,
			ChatId:    req.CallbackQuery.Message.Chat.Id,
			ParseMode: types.TelegramParseModeHTML,
			Text:      textB.String(),
		}, nil

	// Unhandled step
	default:
		// Delete chat
		if repository.TelegramDeleteChat(ctx, chatId) != nil {
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
