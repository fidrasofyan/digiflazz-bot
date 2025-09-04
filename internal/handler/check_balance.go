package handler

import (
	"context"
	"errors"

	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
)

func CheckBalance(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	res, err := service.DigiflazzCheckBalance(ctx)
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

	return &types.TelegramResponse{
		Method:      types.TelegramMethodSendMessage,
		ChatId:      req.Message.Chat.Id,
		ParseMode:   types.TelegramParseModeHTML,
		Text:        util.Sprintf("Saldo Digiflazz: %d", res.Data.Deposit),
		ReplyMarkup: types.DefaultReplyMarkup,
	}, nil
}
