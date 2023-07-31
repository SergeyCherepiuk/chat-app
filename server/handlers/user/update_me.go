package userhandler

import (
	"strings"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
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
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := UpdateUserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.Update(userId, updates); err != nil {
		log.Error(
			"failed to update the user",
			slog.String("err", err.Error()),
			slog.Any("updates", updates),
		)
		return err
	}

	log.Info("user has been updated", slog.Any("updates", updates))
	return c.SendStatus(fiber.StatusOK)
}
