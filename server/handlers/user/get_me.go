package userhandler

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/gofiber/fiber/v2"
)

func (handler UserHandler) GetMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := handler.storage.GetById(userId)
	if err != nil {
		return err
	}

	responseBody := domain.GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	return c.JSON(responseBody)
}
