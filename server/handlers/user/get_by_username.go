package userhandler

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/gofiber/fiber/v2"
)

func (handler UserHandler) GetByUsername(c *fiber.Ctx) error {
	user, err := handler.storage.GetByUsername(c.Params("username"))
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
