package handler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/database"
	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/types"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
)

func Start(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	exists, err := database.Sqlc.IsUserExists(ctx, req.Message.Chat.Id)
	if err != nil {
		return nil, util.NewError(err)
	}

	if exists == 0 {
		_, err := database.Sqlc.CreateUser(ctx, &database.CreateUserParams{
			ID:        req.Message.Chat.Id,
			Username:  &req.Message.Chat.Username,
			FirstName: &req.Message.Chat.FirstName,
			LastName:  &req.Message.Chat.LastName,
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return nil, util.NewError(err)
		}
	}

	return &types.TelegramResponse{
		Method:      types.TelegramMethodSendMessage,
		ChatId:      req.Message.Chat.Id,
		ParseMode:   types.TelegramParseModeHTML,
		Text:        fmt.Sprintf("Welcome to %s! \n\nYour chat ID: <code>%d</code>", config.Cfg.AppName, req.Message.Chat.Id),
		ReplyMarkup: types.DefaultReplyMarkup,
	}, nil
}
