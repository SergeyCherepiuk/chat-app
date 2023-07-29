package userhandler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UpdateUserRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (body UpdateUserRequestBody) ToMap() map[string]any {
	updates := make(map[string]any)

	if strings.TrimSpace(body.FirstName) != "" {
		updates["first_name"] = body.FirstName
	}

	if strings.TrimSpace(body.LastName) != "" {
		updates["last_name"] = body.LastName
	}

	if strings.TrimSpace(body.Username) != "" {
		updates["username"] = body.Username
	}

	return updates
}

func (handler UserHandler) UpdateMe(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)
	
	body := UpdateUserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.Update(userId, updates); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
