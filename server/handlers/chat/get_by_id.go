package chathandler

import (
	"errors"
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) GetById(c *fiber.Ctx) error {
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

	chat, err := handler.storage.GetChatById(uint(chatId))
	if err != nil {
		l.Error("failed to find chat by id", slog.String("err", err.Error()))
		return err
	}

	l.Info("chat has been sent to the user", slog.Any("chat", chat))
	return c.JSON(chat)
}
