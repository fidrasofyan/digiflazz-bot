package repository

import (
	"context"

	"github.com/fidrasofyan/digiflazz-bot/database"
)

func TelegramGetChat(ctx context.Context, id int64) (*database.Chat, error) {
	chat, err := database.Sqlc.GetChat(ctx, id)
	return chat, err
}

type TelegramSetChatParams struct {
	ID      int64
	Command string
	Step    int64
	Data    []byte
}

func TelegramSetChat(ctx context.Context, arg *TelegramSetChatParams) (*database.Chat, error) {
	chatExists, err := database.Sqlc.IsChatExists(ctx, arg.ID)
	if err != nil {
		return nil, err
	}

	if chatExists != 0 {
		chat, err := database.Sqlc.UpdateChat(ctx, &database.UpdateChatParams{
			Command: arg.Command,
			Step:    arg.Step,
			Data:    arg.Data,
			ID:      arg.ID,
		})
		return chat, err
	}

	chat, err := database.Sqlc.CreateChat(ctx, &database.CreateChatParams{
		ID:      arg.ID,
		Command: arg.Command,
		Step:    arg.Step,
		Data:    arg.Data,
	})
	return chat, err
}

type TelegramSetReplyMarkupParams struct {
	ID          int64
	Step        int64
	ReplyMarkup []byte
}

func TelegramSetReplyMarkup(ctx context.Context, params *TelegramSetReplyMarkupParams) error {
	switch params.Step {
	case 1:
		return database.Sqlc.UpdateReplyMarkup1(ctx, &database.UpdateReplyMarkup1Params{
			ID:           params.ID,
			ReplyMarkup1: params.ReplyMarkup,
		})
	case 2:
		return database.Sqlc.UpdateReplyMarkup2(ctx, &database.UpdateReplyMarkup2Params{
			ID:           params.ID,
			ReplyMarkup2: params.ReplyMarkup,
		})
	case 3:
		return database.Sqlc.UpdateReplyMarkup3(ctx, &database.UpdateReplyMarkup3Params{
			ID:           params.ID,
			ReplyMarkup3: params.ReplyMarkup,
		})
	case 4:
		return database.Sqlc.UpdateReplyMarkup4(ctx, &database.UpdateReplyMarkup4Params{
			ID:           params.ID,
			ReplyMarkup4: params.ReplyMarkup,
		})
	default:
		return nil
	}
}

func TelegramDeleteChat(ctx context.Context, id int64) error {
	return database.Sqlc.DeleteChat(ctx, id)
}
