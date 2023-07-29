package userhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	if err := handler.storage.Delete(userId); err != nil {
		l.Error("failed to delete the user", slog.String("err", err.Error()))
		return err
	}

	// TODO: Clear session id in Redis
	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	l.Info("user has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
