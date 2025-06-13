package bot

import (
	"context"
	"log"

	"github.com/akyTheDev/currency-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CmdRegister    = "register"
	CmdDelete      = "delete"
	HelpRegister   = "Register to receive hourly EUR→TRY updates"
	HelpDelete     = "Unregister from receiving updates"
	UnknownCommand = "Unknown command. Use /register or /delete."
)

type BotHandler struct {
	bot           *tgbotapi.BotAPI
	logger        *log.Logger
	userService   *service.UserService
	notifyService *service.NotifyService
	context       context.Context
}

func NewBotHandler(
	context context.Context,
	bot *tgbotapi.BotAPI,
	logger *log.Logger,
	userService *service.UserService,
	notifyService *service.NotifyService,
) *BotHandler {
	return &BotHandler{
		context:       context,
		bot:           bot,
		logger:        logger,
		userService:   userService,
		notifyService: notifyService,
	}
}

func (h *BotHandler) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateConfig)
	go h.startNotify(h.context)
	for {
		select {
		case update := <-updates:
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}
			go h.handleUpdate(update)
		case <-h.context.Done():
			h.logger.Println("BotHandler: stopping due to context cancellation")
			return
		}
	}

}

func (h *BotHandler) handleUpdate(update tgbotapi.Update) {
	msg := update.Message
	chatID := msg.Chat.ID
	cmd := msg.Command()

	h.logger.Printf("Received command: %s from chat_id=%d\n", cmd, chatID)

	switch cmd {
	case CmdRegister:
		h.handleRegister(chatID)
	case CmdDelete:
		h.handleDelete(chatID)
	default:
		h.replyText(chatID, UnknownCommand)
	}
}

func (h *BotHandler) replyText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Printf("replyText: failed to send to chat_id=%d: %v\n", chatID, err)
	}
}
