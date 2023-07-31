package chathandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) GetAll(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chats, err := handler.storage.GetAllChats()
	if err != nil {
		log.Error("failed to get list of chats", slog.String("err", err.Error()))
		return err
	}

	if len(chats) < 1 {
		c.Status(fiber.StatusNoContent)
	} else {
		c.Status(fiber.StatusOK)
	}

	log.Info("list of chats has been sent to the user", slog.Any("chats", chats))
	return c.JSON(chats)
}
