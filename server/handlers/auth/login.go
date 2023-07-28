package authhandler

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	body := domain.LoginRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error("failed to hash the password", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		logger.Logger.Error(
			"request body isn't valid",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	sessionId, userId, err := handler.storage.Login(body.Username, body.Password)
	if err != nil {
		logger.Logger.Error("failed to log in user", slog.String("err", err.Error()))
		return err
	}

	logger.Logger.Info(
		"user has been logged in",
		slog.Uint64("user_id", uint64(userId)),
		slog.Any("session_id", sessionId),
	)
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		HTTPOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return c.SendStatus(fiber.StatusOK)
}
