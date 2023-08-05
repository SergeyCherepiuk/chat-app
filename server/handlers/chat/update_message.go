package chathandler

import (
	chatdomain "github.com/SergeyCherepiuk/chat-app/domain/chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) UpdateMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	messageId := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(messageId)))

	body := chatdomain.UpdateMessageRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error(
			"request body isn't valid",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.chatStorage.UpdateMessage(messageId, body.Message); err != nil {
		log.Error("failed to update the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been updated", slog.Any("updated_message", body.Message))
	return c.SendStatus(fiber.StatusOK)
}
