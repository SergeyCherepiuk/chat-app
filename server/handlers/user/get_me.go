package userhandler

import (
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GetUserResponseBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (handler UserHandler) GetMe(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	user, err := handler.storage.GetById(userId)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to get user by id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return err
	}

	responseBody := GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	logger.LogMessages <- logger.LogMessage{
		Message: "user's info has been sent to the user",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("user", responseBody),
		},
	}
	return c.JSON(responseBody)
}
