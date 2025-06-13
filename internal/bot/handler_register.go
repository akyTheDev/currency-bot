package bot

import (
	"errors"

	"github.com/akyTheDev/currency-bot/internal/domain"
)

func (h *BotHandler) handleRegister(chatID int64) {
	err := h.userService.Register(chatID)

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			h.replyText(chatID, "You are already registered!")
			return
		}
		h.logger.Printf("handleRegister error for chat_id=%d, error: %v\n", chatID, err)
		h.replyText(chatID, "An unexpected error occured. Please try again later.")
		return
	}

	h.replyText(chatID, "✅ You have been registered! You will receive hourly EUR→TRY updates.")
}
