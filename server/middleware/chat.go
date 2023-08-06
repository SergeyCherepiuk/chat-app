package middleware

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/logger"
	chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/exp/slog"
)

type ChatMiddleware struct {
	userStorage userstorage.UserStorage
	chatStorage chatstorage.ChatStorage
}

func NewChatMiddleware(
	userStorage userstorage.UserStorage, chatStorage chatstorage.ChatStorage,
) *ChatMiddleware {
	return &ChatMiddleware{userStorage: userStorage, chatStorage: chatStorage}
}

func (middleware ChatMiddleware) CheckIfCompanionExists() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		username := utils.CopyString(c.Params("username", ""))
		companion, err := middleware.userStorage.GetByUsername(username)
		if err != nil {
			log.Error(
				"failed to find companion",
				slog.String("err", err.Error()),
				slog.String("username", username),
			)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		c.Locals("companion_id", companion.ID)
		return c.Next()
	}
}

func (middleware ChatMiddleware) CheckIfBelongsToChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		userId := c.Locals("user_id").(uint)
		log.With(slog.Uint64("user_id", uint64(userId)))

		companionId := c.Locals("companion_id").(uint)
		log.With(slog.Uint64("companion_id", uint64(companionId)))

		messageId, err := strconv.ParseUint(c.Params("message_id"), 10, 64)
		if err != nil {
			log.Error(
				"failed to parse message id",
				slog.String("err", err.Error()),
				slog.String("message_id", c.Params("message_id")),
			)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		log.With(slog.Uint64("message_id", messageId))

		isBelong, err := middleware.chatStorage.IsMessageBelongToChat(uint(messageId), userId, companionId)
		if err != nil || !isBelong {
			log.Warn("message not belongs to the chat", slog.String("err", err.Error()))
			return c.SendStatus(fiber.StatusBadRequest)
		}

		c.Locals("message_id", uint(messageId))
		return c.Next()
	}
}

func (middleware ChatMiddleware) CheckIfAuthor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		userId := c.Locals("user_id").(uint)
		log.With(slog.Uint64("user_id", uint64(userId)))

		messageId := c.Locals("message_id").(uint)
		log.With(slog.Uint64("message_id", uint64(messageId)))

		isAuthor, err := middleware.chatStorage.IsAuthor(messageId, userId)
		if err != nil {
			log.Error("failed to find out if user is an author", slog.String("err", err.Error()))
			return err
		}

		if !isAuthor {
			log.Warn("user isn't an author of the message")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
