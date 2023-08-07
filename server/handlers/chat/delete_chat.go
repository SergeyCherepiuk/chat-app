package chathandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) DeleteChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	companionId := c.Locals("companion_id").(uint)
	log.With(slog.Uint64("companion_id", uint64(companionId)))

	if err := handler.chatStorage.DeleteAll(userId, companionId); err != nil {
		log.Error("failed to delete the chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
