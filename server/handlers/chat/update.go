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
	userId, _ := c.Locals("user_id").(uint)

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse chat id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("chat_id", c.Params("chat_id")),
			},
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	body := UpdateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse request body",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("chat_id", chatId),
			},
		}
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.UpdateChat(uint(chatId), updates); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to update the chat",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("chat_id", chatId),
				slog.Any("updates", updates),
			},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "chat has been updated",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("chat_id", chatId),
			slog.Any("updates", updates),
		},
	}
	return c.SendStatus(fiber.StatusOK)
}
