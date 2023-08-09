package groupchathandler

import (
	"strconv"

	groupchatdomain "github.com/SergeyCherepiuk/chat-app/domain/group_chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func (handler GroupChatHandler) GetChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		log.Error(
			"failed to parse group chat id",
			slog.String("err", err.Error()),
			slog.String("chat_id", c.Params("chat_id")),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	chat, err := handler.storage.GetChat(uint(chatId))
	if err != nil {
		slog.Error("failed to get the group chat info", slog.String("err", err.Error()))
		return err
	}

	responseBody := groupchatdomain.GetGroupChatResponseBody{
		Name: chat.Name,
	}
	log.Info("group chat info has been sent to the user")
	return c.JSON(responseBody)
}
