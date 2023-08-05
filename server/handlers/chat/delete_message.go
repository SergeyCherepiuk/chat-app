package chathandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) DeleteMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	message_id := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(message_id)))

	if err := handler.chatStorage.DeleteMessage(message_id); err != nil {
		log.Error("failed to delete the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
