package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akyTheDev/currency-bot/internal/bot"
	"github.com/akyTheDev/currency-bot/internal/config"
	"github.com/akyTheDev/currency-bot/internal/fetcher"
	"github.com/akyTheDev/currency-bot/internal/repository"
	"github.com/akyTheDev/currency-bot/internal/service"
	"github.com/akyTheDev/currency-bot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Config couldn't be loaded: %v", err)
	}

	db, err := storage.OpenDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Couldn't connected to db: %v", err)
	}

	// Fetcher
	fetcher := fetcher.NewTCMBClient(fetcher.TcmbUrl, 60)

	// Repositories
	userRepository := repository.NewPostgresUserRepository(db)

	// Services
	userService := service.NewUserService(userRepository, logger)
	notifyService := service.NewNotifyService(logger, userRepository, fetcher)

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Fatalf("Couldn't create a bot")
	}

	commands := []tgbotapi.BotCommand{
		{Command: bot.CmdRegister, Description: bot.HelpRegister},
		{Command: bot.CmdDelete, Description: bot.HelpDelete},
	}

	setCmdConfig := tgbotapi.NewSetMyCommands(commands...)
	if _, err := botAPI.Request(setCmdConfig); err != nil {
		logger.Fatalf("could not set bot commands: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	handler := bot.NewBotHandler(ctx, botAPI, logger, userService, notifyService)
	go handler.Start()

	logger.Println("Bot is running...")
	<-ctx.Done()
	logger.Println("Shutting down...")
	time.Sleep(1 * time.Second)
}
