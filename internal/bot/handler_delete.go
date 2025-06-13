package bot

import (
	"errors"

	"github.com/akyTheDev/currency-bot/internal/domain"
)

func (h *BotHandler) handleDelete(chatID int64) {
	err := h.userService.Delete(chatID)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.replyText(chatID, "You are not registered!")
			return
		}
		h.logger.Printf("handleDelete error for chat_id=%d, error: %v\n", chatID, err)
		h.replyText(chatID, "An unexpected error occured. Please try again later.")
		return
	}

	h.replyText(chatID, "🗑️ You have been unregistered. You will no longer receive updates.")
}
