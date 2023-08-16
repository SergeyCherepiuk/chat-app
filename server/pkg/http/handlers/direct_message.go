package handlers

import (
	"context"
	"fmt"
	"math"
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

type DirectMessageHandler struct {
	directMessageService     domain.DirectMessageService
	connectionManagerService *connection.ConnectionManagerService[[2]uint]
	messageSenderService     *messaging.MessageSenderService[domain.DirectMessage]
	userService              domain.UserService
}

func NewDirectMessageHandler(
	directMessageService domain.DirectMessageService,
	userService domain.UserService,
) *DirectMessageHandler {
	return &DirectMessageHandler{
		directMessageService:     directMessageService,
		connectionManagerService: connection.NewConnectionManager[[2]uint](),
		messageSenderService:     messaging.NewMessageSenderService[domain.DirectMessage](),
		userService:              userService,
	}
}

type Pair struct {
	First  uint
	Second uint
}

func (pair Pair) GetKey() [2]uint {
	if pair.First > pair.Second {
		return [2]uint{pair.Second, pair.First}
	}
	return [2]uint{pair.First, pair.Second}
}

func (handler DirectMessageHandler) EnterChat(c *websocket.Conn) {
	defer c.Close()

	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	companionId := c.Locals("companion_id").(uint)
	log.With(slog.Uint64("companion_id", uint64(companionId)))

	key := Pair{First: userId, Second: companionId}.GetKey()
	go handler.connectionManagerService.Connect(key, c)
	log.Info("user has been connected to the chat", slog.Any("key", key))
	defer func() {
		go handler.connectionManagerService.Disconnect(key, c)
	}()

	history, err := handler.directMessageService.GetHistory(userId, companionId, math.MaxInt64)
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
			log.Error("failed to read the message", slog.String("err", err.Error()))
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

		message := domain.DirectMessage{
			Message:  body.Message,
			From:     userId,
			To:       companionId,
			IsEdited: false,
		}
		if err := handler.directMessageService.Create(&message); err != nil {
			log.Error(
				"failed to store the message",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}

		go handler.messageSenderService.Send(
			context.WithValue(context.Background(), logger.LogContextKey, log),
			[]domain.DirectMessage{message},
			handler.connectionManagerService.GetConnections(key)...,
		)
		log.Info("user has sent the message", slog.Any("message", message))
	}
}

func (handler DirectMessageHandler) GetHistory(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	companionId := c.Locals("companion_id").(uint)
	log.With(slog.Uint64("companion_id", uint64(companionId)))

	fromId, err := strconv.ParseUint(c.Query("from_id"), 10, 64)
	if c.Query("from_id") != "" && err != nil {
		log.Error("failed to parse 'from_id' message id", slog.String("err", err.Error()))
		return err
	}
	if c.Query("from_id") == "" {
		fromId = math.MaxInt64
	}

	history, err := handler.directMessageService.GetHistory(userId, companionId, uint(fromId))
	if err != nil {
		log.Error("failed to get chat history", slog.String("err", err.Error()))
		return err
	}

	defer log.Info("chat history has been sent")
	if len(history) != settings.CHAT_HISTORY_BLOCK_SIZE {
		return c.JSON(validation.GetHistoryResponseBody{History: history})
	}

	next := fmt.Sprintf(
		"http://localhost:8001/api/chat/%s/history?from_id=%d",
		c.Params("username"),
		history[len(history)-1].ID-1,
	)
	return c.JSON(validation.GetHistoryWithNextResponseBody{
		History: history,
		Next:    string(next),
	})
}

func (handler DirectMessageHandler) UpdateMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	messageId := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(messageId)))

	body := validation.UpdateMessageRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error(
			"request body isn't valid",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.directMessageService.Update(messageId, body.Message); err != nil {
		log.Error("failed to update the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been updated", slog.Any("updated_message", body.Message))
	return c.SendStatus(fiber.StatusOK)
}

func (handler DirectMessageHandler) DeleteMessage(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	message_id := c.Locals("message_id").(uint)
	log.With(slog.Uint64("message_id", uint64(message_id)))

	if err := handler.directMessageService.Delete(message_id); err != nil {
		log.Error("failed to delete the message", slog.String("err", err.Error()))
		return err
	}

	log.Info("message has been deleted")
	return c.SendStatus(fiber.StatusOK)
}

func (handler DirectMessageHandler) DeleteChat(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	companionId := c.Locals("companion_id").(uint)
	log.With(slog.Uint64("companion_id", uint64(companionId)))

	if err := handler.directMessageService.DeleteAll(userId, companionId); err != nil {
		log.Error("failed to delete the chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
