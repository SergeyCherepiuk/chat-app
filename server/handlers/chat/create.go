package chathandler

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		logger.Logger.Error("failed to parse user id", slog.Any("user_id", c.Locals("user_id")))
		return errors.New("failed to parse user id")
	}
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	body := domain.CreateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		l.Error(
			"failed to parse request body",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	chat := models.Chat{Name: body.Name}
	if err := handler.storage.CreateChat(&chat); err != nil {
		l.Error("failed to create new chat", slog.String("err", err.Error()))
		return err
	}

	l.Info("new chat has been created", slog.Any("chat", chat))
	return c.SendStatus(fiber.StatusOK)
}
