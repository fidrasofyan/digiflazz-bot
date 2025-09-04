package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/fidrasofyan/digiflazz-bot/cmd"
	"github.com/fidrasofyan/digiflazz-bot/database"
	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/fidrasofyan/digiflazz-bot/internal/job"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
	"github.com/fidrasofyan/digiflazz-bot/internal/util"
	"github.com/gofiber/fiber/v2"
)

var AppVersion = "n/a"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: digiflazz-bot <command>. See 'digiflazz-bot help' for more info.")
		os.Exit(1)
	}

	mainCtx, cancel := context.WithCancel(context.Background())

	// Load config
	config.MustLoadConfig()

	// Setup signal catching
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh,
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // stop
		syscall.SIGQUIT, // Ctrl+\
		syscall.SIGHUP,  // terminal hangup
	)
	errCh := make(chan error, 1)

	var httpServer *fiber.App

	switch os.Args[1] {
	case "start":
		go func() {
			log.Printf(
				"Environment: %s - Runtime: %s - App: %s - TZ: %s\n",
				config.Cfg.AppEnv,
				runtime.Version(),
				AppVersion,
				os.Getenv("TZ"),
			)

			// Load database
			database.MustLoadDatabase(mainCtx)

			// Set webhook
			if config.Cfg.AppEnv == "production" {
				err := cmd.SetTelegramWebhookAndCommands(mainCtx)
				if err != nil {
					errCh <- err
					return
				}
			}

			// Start HTTP server
			httpServer = cmd.MustStartHTTPServer()
		}()

	case "set-telegram-webhook":
		go func() {
			err := cmd.SetTelegramWebhookAndCommands(mainCtx)
			if err != nil {
				errCh <- err
			}
			quitCh <- syscall.SIGQUIT
		}()

	case "populate-products":
		go func() {
			// Load database
			database.MustLoadDatabase(mainCtx)

			// Populate products
			err := job.PopulateProducts(mainCtx)
			if err != nil {
				errCh <- err
				return
			}
			quitCh <- syscall.SIGQUIT
		}()

	case "digiflazz-sign":
		if len(os.Args) < 3 {
			errCh <- errors.New("missing second argument")
			return
		}
		sign := service.DigiflazzSign(os.Args[2])
		fmt.Println("Sign:", sign)
		quitCh <- syscall.SIGQUIT

	case "generate-secret":
		if len(os.Args) < 3 {
			errCh <- errors.New("missing second argument")
			return
		}
		// It must be numeric
		length, err := strconv.Atoi(os.Args[2])
		if err != nil {
			errCh <- err
			return
		}
		secret, err := util.GenerateSecretToken(length)
		if err != nil {
			errCh <- err
			return
		}
		fmt.Println("Secret:", secret)
		quitCh <- syscall.SIGQUIT

	case "help":
		help := []string{
			"start                          Start the bot",
			"set-telegram-webhook           Set Telegram webhook and commands",
			"populate-products              Populate products",
			"digiflazz-sign <string>        Generate Digiflazz sign",
			"generate-secret <int>          Generate secret token",
			"help                           Show this help",
		}
		fmt.Println("Commands:")
		for _, h := range help {
			fmt.Printf("  %s\n", h)
		}
		quitCh <- syscall.SIGQUIT

	default:
		fmt.Printf("Unknown command '%s'. Use 'help' to see available commands.\n", os.Args[1])
		os.Exit(1)
	}

	// Wait for signal
	select {
	case sig := <-quitCh:
		if os.Args[1] == "start" {
			log.Printf("Signal caught: %s", sig)
		}
	case err := <-errCh:
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}

	// Cancel context
	cancel()
	<-mainCtx.Done()

	// Stop HTTP server
	if httpServer != nil {
		log.Println("Stopping HTTP server...")
		err := httpServer.Shutdown()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	// Close database
	if database.DBConn != nil {
		log.Println("Closing database...")
		database.DBConn.Close()
	}

	if os.Args[1] == "start" {
		log.Println("Goodbye!")
	}
}
