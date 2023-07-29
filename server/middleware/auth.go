package middleware

import (
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userId, err := middleware.storage.Check(sessionId)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user_id", userId)
		return c.Next()
	}
}
