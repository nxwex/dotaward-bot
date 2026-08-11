package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/telebot.v3"
)

func (h *Handler) Clear(c telebot.Context) error {
	if !isAdmin(c) {
		return c.Reply("Недостаточно прав")
	}

	contexts, err := h.contexts.GetMatchContexts(c.Chat().ID)
	if err != nil {
		log.Printf("GetMatchContexts error: %v", err)
		return c.Reply("Ошибка получения сообщений")
	}

	deleted := 0

	for _, a := range contexts {
		if err := c.Bot().Delete(&telebot.Message{
			ID:   a.MessageID,
			Chat: c.Chat(),
		}); err != nil {
			log.Printf(
				"Delete message error: chatID=%d messageID=%d err=%v",
				c.Chat().ID,
				a.MessageID,
				err,
			)
			continue
		}

		deleted++
	}

	if err := h.contexts.DeleteMatchContexts(c.Chat().ID); err != nil {
		log.Printf("DeleteMatchContexts error: %v", err)
		return c.Reply("Сообщения удалены, но контекст не очищен")
	}

	return c.Reply(fmt.Sprintf("Удалено сообщений: %d", deleted))
}

func isAdmin(c telebot.Context) bool {
	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		return false
	}

	adminChatID, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		return false
	}

	return c.Sender() != nil &&
		c.Chat() != nil &&
		c.Sender().ID == adminID &&
		c.Chat().ID == adminChatID
}
