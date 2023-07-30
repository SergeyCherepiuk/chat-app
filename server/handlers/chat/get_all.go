package chathandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) GetAll(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	chats, err := handler.storage.GetAllChats()
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to get list of chats",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return err
	}

	if len(chats) < 1 {
		c.Status(fiber.StatusNoContent)
	} else {
		c.Status(fiber.StatusOK)
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "list of chats has been sent to the user",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("chats", chats),
		},
	}
	return c.JSON(chats)
}
