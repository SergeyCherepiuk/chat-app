package userhandler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	if err := handler.storage.Delete(userId); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
