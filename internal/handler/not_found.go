package handler

import (
	"context"

	"github.com/fidrasofyan/digiflazz-bot/database/repository"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
)

func NotFound(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	// Is it callback query?
	if req.CallbackQuery != nil {
		// Delete chat
		err := repository.TelegramDeleteChat(ctx, req.CallbackQuery.From.Id)
		if err != nil {
			return nil, util.NewError(err)
		}

		// Answer callback query
		err = service.TelegramAnswerCallbackQuery(ctx, &service.TelegramAnswerCallbackQueryParams{
			CallbackQueryId: req.CallbackQuery.Id,
		})
		if err != nil {
			return nil, util.NewError(err)
		}

		return &types.TelegramResponse{
			Method:    types.TelegramMethodEditMessageText,
			MessageId: req.CallbackQuery.Message.MessageId,
			ChatId:    req.CallbackQuery.Message.Chat.Id,
			ParseMode: types.TelegramParseModeHTML,
			Text:      "<i>Sesi tidak valid</i>",
		}, nil
	}

	return &types.TelegramResponse{
		Method:      types.TelegramMethodSendMessage,
		ChatId:      req.Message.Chat.Id,
		ParseMode:   types.TelegramParseModeHTML,
		Text:        "<i>Perintah tidak dikenali</i>",
		ReplyMarkup: types.DefaultReplyMarkup,
	}, nil
}
