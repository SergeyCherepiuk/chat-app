package middleware

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	groupchatstorage "github.com/SergeyCherepiuk/chat-app/storage/group_chat"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GroupChatMiddleware struct {
	storage groupchatstorage.GroupChatStorage
}

func NewGroupChatMiddleware(storage groupchatstorage.GroupChatStorage) *GroupChatMiddleware {
	return &GroupChatMiddleware{storage: storage}
}

func (middleware GroupChatMiddleware) CheckIfAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			return c.SendStatus(fiber.StatusBadRequest)
		}

		isAdmin, err := middleware.storage.IsAdmin(uint(chatId), userId)
		if err != nil {
			log.Error("failed to find out if user is admin", slog.String("err", err.Error()))
			return err
		}

		if !isAdmin {
			log.Warn("user isn't a creator of the group chat")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("chat_id", uint(chatId))
		return c.Next()
	}
}
