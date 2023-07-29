package userhandler

import (
	"github.com/gofiber/fiber/v2"
)

type GetMeResponseBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (handler UserHandler) GetByUsername(c *fiber.Ctx) error {
	user, err := handler.storage.GetByUsername(c.Params("username"))
	if err != nil {
		return err
	}

	responseBody := GetMeResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	return c.JSON(responseBody)
}
