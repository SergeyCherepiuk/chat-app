package middleware

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type AuthMiddleware struct {
	authService domain.AuthService
}

func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (middleware AuthMiddleware) CheckIfAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
		if err != nil {
			log.Error(
				"failed to update the user",
				slog.String("err", err.Error()),
				slog.String("session_id", c.Cookies("session_id", "")),
			)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userId, err := middleware.authService.Check(sessionId)
		if err != nil {
			log.Error(
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
