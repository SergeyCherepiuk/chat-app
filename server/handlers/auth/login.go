package authhandler

import (
	"time"

	authdomain "github.com/SergeyCherepiuk/chat-app/domain/auth"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.Logger{}

	body := authdomain.LoginRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse the body", slog.String("err", err.Error()))
		return err
	}

	sessionId, userId, err := handler.storage.Login(body.Username, body.Password)
	if err != nil {
		log.Error("failed to log in user", slog.String("err", err.Error()))
		return err
	}

	log.Info(
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
