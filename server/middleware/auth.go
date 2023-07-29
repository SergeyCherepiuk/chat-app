package middleware

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type AuthMiddleware struct {
	storage authstorage.AuthStorage
}

func NewAuthMiddleware(storage authstorage.AuthStorage) *AuthMiddleware {
	return &AuthMiddleware{storage: storage}
}

func (middleware AuthMiddleware) CheckIfAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
		if err != nil {
			logger.Logger.Error(
				"failed to parse UUID",
				slog.String("err", err.Error()),
				slog.String("session_id", c.Cookies("session_id", "")),
			)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userId, err := middleware.storage.Check(sessionId)
		if err != nil {
			logger.Logger.Error(
				"failed to find a session",
				slog.String("err", err.Error()),
				slog.Any("session_id", sessionId),
			)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user_id", userId)
		return c.Next()
	}
}
