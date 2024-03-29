package handlers

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/connection"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/SergeyCherepiuk/chat-app/pkg/messaging"
	"github.com/SergeyCherepiuk/chat-app/pkg/settings"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GroupChatHandler struct {
	groupChatService         domain.GroupChatService
	connectionManagerService *connection.ConnectionManagerService[uint]
	messageSenderService     *messaging.MessageSenderService[domain.GroupMessage]
}

func NewGroupChatHandler(groupChatService domain.GroupChatService) *GroupChatHandler {
	return &GroupChatHandler{
		groupChatService:         groupChatService,
		connectionManagerService: connection.NewConnectionManager[uint](),
		messageSenderService:     &messaging.MessageSenderService[domain.GroupMessage]{},
	}
}

func (handler GroupChatHandler) EnterChat(c *websocket.Conn) {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	go handler.connectionManagerService.Connect(chatId, c)
	log.Info("user has been connected to the chat")
	defer func() {
		go handler.connectionManagerService.Disconnect(chatId, c)
	}()

	history, err := handler.groupChatService.GetHistory(chatId, math.MaxInt64)
	if err != nil {
		log.Error("failed to get chat history", slog.String("err", err.Error()))
		return
	}

	historyContext, historyCancel := context.WithCancel(context.Background())
	historyContext = context.WithValue(historyContext, logger.LogContextKey, log)
	defer historyCancel()
	go handler.messageSenderService.Send(historyContext, history, c)

	for {
		_, text, err := c.ReadMessage()
		if websocket.IsCloseError(err, 1000, 1005) {
			log.Info("user has been disconnected")
			return
		}
		if err != nil {
			slog.Error("failed to read message", slog.String("err", err.Error()))
			return
		}

		body := validation.CreateMessageBody{Message: string(text)}
		if err := body.Validate(); err != nil {
			log.Error(
				"body isn't valid",
				slog.String("err", err.Error()),
				slog.Any("body", body),
			)
			continue
		}

		message := domain.GroupMessage{
			Message:  body.Message,
			UserID:   userId,
			ChatID:   chatId,
			IsEdited: false,
		}
		if err := handler.groupChatService.CreateMessage(&message); err != nil {
			log.Error(
				"failed to store the message",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}

		go handler.messageSenderService.Send(
			context.WithValue(context.Background(), logger.LogContextKey, log),
			[]domain.GroupMessage{message},
			handler.connectionManagerService.GetConnections(chatId)...,
		)
		log.Info("user has sent the message", slog.Any("message", message))
	}
}

func (handler GroupChatHandler) GetHistory(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	fromId, err := strconv.ParseUint(c.Query("from_id"), 10, 64)
	if c.Query("from_id") != "" && err != nil {
		log.Error("failed to parse 'from_id' message id", slog.String("err", err.Error()))
		return err
	}
	if c.Query("from_id") == "" {
		fromId = math.MaxInt64
	}

	history, err := handler.groupChatService.GetHistory(chatId, uint(fromId))
	if err != nil {
		log.Error("failed to get chat history", slog.String("err", err.Error()))
		return err
	}

	defer log.Info("chat history has been sent")
	if len(history) < settings.CHAT_HISTORY_BLOCK_SIZE {
		return c.JSON(validation.GetHistoryResponseBody[domain.GroupMessage]{
			History: history,
		})
	}

	next := fmt.Sprintf(
		"http://localhost:%s/api/group-chat/%d/history?from_id=%d",
		os.Getenv("SERVER_PORT"),
		chatId,
		history[len(history)-1].ID-1,
	)
	return c.JSON(validation.GetHistoryWithNextResponseBody[domain.GroupMessage]{
		History: history,
		Next:    next,
	})
}

func (handler GroupChatHandler) GetChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	chat, err := handler.groupChatService.GetChat(uint(chatId))
	if err != nil {
		slog.Error("failed to get the group chat info", slog.String("err", err.Error()))
		return err
	}

	responseBody := validation.GetGroupChatResponseBody{
		Name: chat.Name,
	}
	log.Info("group chat info has been sent to the user")
	return c.JSON(responseBody)
}

func (handler GroupChatHandler) CreateChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := validation.CreateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error("invalid request body", slog.String("err", err.Error()))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	chat := domain.GroupChat{
		Name:      body.Name,
		CreatorID: userId,
	}
	if err := handler.groupChatService.CreateChat(&chat); err != nil {
		log.Error("failed to store the group chat", slog.String("err", err.Error()))
	}

	log.Info("group chat has been stored", slog.Any("chat", chat))
	return nil
}

func (handler GroupChatHandler) UpdateChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	body := validation.UpdateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	log.With(slog.Any("updates", updates))

	if err := handler.groupChatService.UpdateChat(chatId, updates); err != nil {
		log.Error("failed to update the group chat", slog.String("err", err.Error()))
		return err
	}

	slog.Info("group chat has been updated")
	return c.SendStatus(fiber.StatusOK)
}

func (handler GroupChatHandler) UpdateMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	messageId := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(messageId)))

	body := validation.UpdateMessageRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := body.Validate(); err != nil {
		log.Error(
			"invalid request body",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.groupChatService.UpdateMessage(messageId, body.Message); err != nil {
		log.Error("failed to update the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been updated")
	return c.SendStatus(fiber.StatusOK)
}

func (handler GroupChatHandler) DeleteChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	if err := handler.groupChatService.DeleteChat(chatId); err != nil {
		log.Error("failed to delete the group chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("group chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}

func (handler GroupChatHandler) DeleteMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	messageId := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(messageId)))

	if err := handler.groupChatService.DeleteMessage(messageId); err != nil {
		log.Error("failed to delete the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
