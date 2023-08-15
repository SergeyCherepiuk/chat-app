package middleware

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GroupChatMiddleware struct {
	groupChatService domain.GroupChatService
}

func NewGroupChatMiddleware(groupChatService domain.GroupChatService) *GroupChatMiddleware {
	return &GroupChatMiddleware{groupChatService: groupChatService}
}

func (middleware GroupChatMiddleware) CheckIfGroupChatExists() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
		if err != nil {
			log.Error(
				"failed to parse group chat id",
				slog.String("err", err.Error()),
				slog.String("chat_id", c.Params("chat_id")),
			)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		c.Locals("chat_id", uint(chatId))
		return c.Next()
	}
}

func (middleware GroupChatMiddleware) CheckIfAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		userId := c.Locals("user_id").(uint)
		log.With(slog.Uint64("user_id", uint64(userId)))

		chatId := c.Locals("chat_id").(uint)
		log.With(slog.Uint64("chat_id", uint64(chatId)))

		isAdmin, err := middleware.groupChatService.IsAdminOfChat(chatId, userId)
		if err != nil {
			log.Error("failed to find out if user is admin", slog.String("err", err.Error()))
			return err
		}

		if !isAdmin {
			log.Warn("user isn't a creator of the group chat")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func (middleware GroupChatMiddleware) CheckIfMessageBelongsToChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		chatId := c.Locals("chat_id").(uint)
		log.With(slog.Uint64("chat_id", uint64(chatId)))

		messageId, err := strconv.ParseUint(c.Params("message_id"), 10, 64)
		if err != nil {
			log.Error(
				"failed to parse message id",
				slog.String("err", err.Error()),
				slog.Uint64("message_id", messageId),
			)
			return c.Status(fiber.StatusBadRequest).SendString(err.Error()) 
		}
		
		belongs, err := middleware.groupChatService.IsMessageBelongsToChat(uint(messageId), chatId)
		if err != nil || !belongs {
			log.Warn("message not belongs to the chat", slog.String("err", err.Error()))
			return c.SendStatus(fiber.StatusBadRequest)
		}
		
		c.Locals("message_id", uint(messageId))
		return c.Next()
	}
}

func (middleware GroupChatMiddleware) CheckIfAuthorOfMessage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.Logger{}

		userId := c.Locals("user_id").(uint)
		log.With(slog.Uint64("user_id", uint64(userId)))

		messageId := c.Locals("message_id").(uint)
		log.With(slog.Uint64("message_id", uint64(messageId)))

		isAuthor, err := middleware.groupChatService.IsAuthorOfMessage(messageId, userId)
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