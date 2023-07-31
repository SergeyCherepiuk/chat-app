package chathandler

import (
	"strconv"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type UpdateChatRequestBody struct {
	Name string `json:"name"`
}

func (body UpdateChatRequestBody) ToMap() map[string]any {
	updates := make(map[string]any)

	if strings.TrimSpace(body.Name) != "" {
		updates["name"] = body.Name
	}

	return updates
}

func (handler ChatHandler) Update(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		log.Error(
			"chat has been sent to the user",
			slog.String("err", err.Error()),
			slog.Any("chat_id", c.Params("chat_id")),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	log.With(slog.Uint64("chat_id", chatId))

	body := UpdateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.UpdateChat(uint(chatId), updates); err != nil {
		log.Error(
			"failed to update the chat",
			slog.String("err", err.Error()),
			slog.Any("updates", updates),
		)
		return err
	}

	log.Info("chat has been updated", slog.Any("updates", updates))
	return c.SendStatus(fiber.StatusOK)
}
