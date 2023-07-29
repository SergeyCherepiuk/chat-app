package chathandler

import (
	"errors"
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
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		logger.Logger.Error("failed to parse user id", slog.Any("user_id", c.Locals("user_id")))
		return errors.New("failed to parse user id")
	}
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		l.Error(
			"failed to parse chat id",
			slog.String("err", err.Error()),
			slog.Any("chat_id", c.Params("chat_id")),
		)
		return err
	}
	l = l.With(slog.Uint64("chat_id", chatId))

	body := UpdateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		l.Error(
			"failed to parse request body",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.UpdateChat(uint(chatId), updates); err != nil {
		l.Error(
			"failed to update the chat",
			slog.String("err", err.Error()),
			slog.Any("updates", updates),
		)
		return err
	}

	l.Info("chat has been updated", slog.Any("updates", updates))
	return c.SendStatus(fiber.StatusOK)
}
