package authhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "invalid session id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Any("session_id", sessionId),
			},
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.storage.Logout(sessionId); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to log out user",
			Level:   slog.LevelError,
			Attrs:   []slog.Attr{slog.String("err", err.Error())},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "user has been logged out",
		Level:   slog.LevelInfo,
		Attrs:   []slog.Attr{slog.Any("session_id", sessionId)},
	}
	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
