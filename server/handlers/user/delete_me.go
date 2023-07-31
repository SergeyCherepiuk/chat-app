package userhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	if err := handler.storage.Delete(userId); err != nil {
		log.Error("failed to delete the user", slog.String("err", err.Error()))
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	log.Info("user has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
