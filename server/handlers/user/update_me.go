package userhandler

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/gofiber/fiber/v2"
)

func (handler UserHandler) UpdateMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	body := domain.UpdateUserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.Update(userId, updates); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
