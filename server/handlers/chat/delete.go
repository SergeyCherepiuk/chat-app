package chathandler

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) Delete(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse chat id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Any("chat_id", c.Params("chat_id")),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.storage.DeleteChat(uint(chatId)); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to delete the chat",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Uint64("chat_id", chatId),
			},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "chat has been deleted",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Uint64("chat_id", chatId),
		},
	}
	return c.SendStatus(fiber.StatusOK)
}
