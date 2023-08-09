package groupchathandler

import (
	groupchatdomain "github.com/SergeyCherepiuk/chat-app/domain/group_chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler GroupChatHandler) Create(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := groupchatdomain.CreateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error("invalid request body", slog.String("err", err.Error()))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	chat := models.GroupChat{
		Name: body.Name,
		CreatorID: userId,
	}
	if err := handler.storage.Create(&chat); err != nil {
		log.Error("failed to store the group chat", slog.String("err", err.Error()))
	}

	log.Info("group chat has been stored", slog.Any("chat", chat))
	return nil
}
