package userhandler

import (
	userdomain "github.com/SergeyCherepiuk/chat-app/domain/user"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler UserHandler) UpdateMe(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := userdomain.UpdateUserRequestBody{}
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
