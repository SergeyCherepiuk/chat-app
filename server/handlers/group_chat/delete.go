package groupchathandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler GroupChatHandler) Delete(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))
	
	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	if err := handler.storage.Delete(chatId); err != nil {
		log.Error("failed to delete the group chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("group chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}