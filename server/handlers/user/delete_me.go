package userhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	if err := handler.storage.Delete(userId); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to delete the user",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return err
	}

	// TODO: Clear session id in Redis
	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	logger.LogMessages <- logger.LogMessage{
		Message: "user has been deleted",
		Level:   slog.LevelInfo,
		Attrs:   []slog.Attr{slog.Uint64("user_id", uint64(userId))},
	}
	return c.SendStatus(fiber.StatusOK)
}
