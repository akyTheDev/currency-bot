package bot

import (
	"context"
	"fmt"
	"time"
)

const interval = time.Hour

func (h *BotHandler) startNotify(ctx context.Context) {
	ticker := time.NewTicker(interval)

	go h.notify()

	for {
		select {
		case <-ticker.C:
			h.notify()
		case <-ctx.Done():
			h.logger.Println("NotifyHandler: context canceled; stopping notifications")
			ticker.Stop()
			return
		}
	}
}

func (h *BotHandler) notify() {
	h.logger.Println("NOTIFY TRIGGERED.")
	ids, rate, err := h.notifyService.GetUsersAndCurrencyRate()
	if err != nil {
		h.logger.Printf("NotifyHandler: failed to get users or rate: %v\n", err)
		return
	}

	if len(ids) == 0 {
		h.logger.Println("NotifyHandler: no subscribers; skipping send")
		return
	}

	text := fmt.Sprintf("EUR→TRY Selling: %.4f Buying: %.4f (at %s)", rate.Selling, rate.Buying, time.Now().Format("15:04"))
	for _, chatID := range ids {
		h.replyText(chatID, text)
	}
	h.logger.Printf("%d users have been notified.\n", len(ids))
}
