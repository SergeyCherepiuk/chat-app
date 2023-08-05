package userhandler

import (
	userdomain "github.com/SergeyCherepiuk/chat-app/domain/user"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler UserHandler) GetByUsername(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	user, err := handler.storage.GetByUsername(c.Params("username"))
	if err != nil {
		log.Error(
			"failed to get user by username",
			slog.String("err", err.Error()),
			slog.String("username", c.Params("username")),
		)
		return err
	}

	responseBody := userdomain.GetUserResponseBody{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Description: user.Description,
	}
	log.Info("user's info has been sent to the user", slog.Any("user", responseBody))
	return c.JSON(responseBody)
}
