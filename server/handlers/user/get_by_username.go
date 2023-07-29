package userhandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GetMeResponseBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (handler UserHandler) GetByUsername(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	user, err := handler.storage.GetByUsername(c.Params("username"))
	if err != nil {
		l.Error(
			"failed to get user by username",
			slog.String("err", err.Error()),
			slog.String("username", c.Params("username")),
		)
		return err
	}

	responseBody := GetMeResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	l.Info(
		"user's info has been sent to the user",
		slog.Any("user", responseBody),
		slog.String("username", c.Params("username")),
	)
	return c.JSON(responseBody)
}
