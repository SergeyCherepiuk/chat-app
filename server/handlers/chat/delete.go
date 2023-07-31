package chathandler

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler ChatHandler) Delete(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		log.Error(
			"failed to parse chat id",
			slog.String("err", err.Error()),
			slog.Any("chat_id", c.Params("chat_id")),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	log.With(slog.Uint64("chat_id", chatId))

	if err := handler.storage.DeleteChat(uint(chatId)); err != nil {
		log.Error("failed to delete the chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
