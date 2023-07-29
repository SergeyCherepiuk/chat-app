package chathandler

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) Delete(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)
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

	if err := handler.storage.DeleteChat(uint(chatId)); err != nil {
		l.Error("failed to delete the chat", slog.String("err", err.Error()))
		return err
	}

	l.Info("chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
