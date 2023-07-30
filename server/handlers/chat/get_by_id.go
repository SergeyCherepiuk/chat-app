package chathandler

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) GetById(c *fiber.Ctx) error {
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

	chat, err := handler.storage.GetChatById(uint(chatId))
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to find chat by id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("chat_id", chatId),
			},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "chat has been sent to the user",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("chat_id", chatId),
			slog.Any("chat", chat),
		},
	}
	return c.JSON(chat)
}
