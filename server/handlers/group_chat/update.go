package groupchathandler

import (
	groupchatdomain "github.com/SergeyCherepiuk/chat-app/domain/groupchat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler GroupChatHandler) Update(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	body := groupchatdomain.UpdateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	log.With(slog.Any("updates", updates))

	if err := handler.storage.Update(chatId, updates); err != nil {
		log.Error("failed to update the group chat", slog.String("err", err.Error()))
		return err
	}

	slog.Info("group chat has been updated")
	return c.SendStatus(fiber.StatusOK)
}
