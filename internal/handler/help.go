package handler

import (
	"context"
	"strings"

	"github.com/fidrasofyan/digiflazz-bot/internal/types"
)

func Help(ctx context.Context, req *types.TelegramUpdate) (*types.TelegramResponse, error) {
	var textB strings.Builder
	textB.WriteString("Format transaksi: kode_produk nomor_tujuan\n")
	textB.WriteString("Contoh: <code>IG100 085808580858</code>")

	return &types.TelegramResponse{
		Method:      types.TelegramMethodSendMessage,
		ChatId:      req.Message.Chat.Id,
		ParseMode:   types.TelegramParseModeHTML,
		Text:        textB.String(),
		ReplyMarkup: types.DefaultReplyMarkup,
	}, nil
}
