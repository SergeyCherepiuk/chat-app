package userhandler

import (
	"github.com/gofiber/fiber/v2"
)

type GetUserResponseBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (handler UserHandler) GetMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := handler.storage.GetById(userId)
	if err != nil {
		return err
	}

	responseBody := GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	return c.JSON(responseBody)
}
