package authhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	log := logger.Logger{}

	sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
	if err != nil {
		log.Error(
			"invalid session id",
			slog.String("err", err.Error()),
			slog.Any("session_id", sessionId),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.storage.Logout(sessionId); err != nil {
		log.Error("failed to log out user", slog.String("err", err.Error()))
		return err
	}

	log.Info("user has been logged out", slog.Any("session_id", sessionId))
	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
