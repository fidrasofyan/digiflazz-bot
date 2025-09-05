package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/fidrasofyan/digiflazz-bot/internal/service"
)

func SetTelegramWebhookAndCommands(ctx context.Context) error {
	log.Println("Setting webhook and commands...")
	// Set webhook
	err := service.TelegramSetWebhook(ctx)
	if err != nil {
		return fmt.Errorf("setting webhook: %v", err)
	}

	// Set commands
	commands := []service.TelegramCommand{
		{
			Command:     "/start",
			Description: "Start the bot",
		},
		{
			Command:     "/help",
			Description: "Show help message",
		},
	}
	err = service.TelegramSetMyCommands(ctx, commands)
	if err != nil {
		return fmt.Errorf("setting commands: %v", err)
	}

	log.Println("DONE: setting webhook and commands")
	return nil
}
