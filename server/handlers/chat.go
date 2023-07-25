package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type ChatHandler struct {
	storage *storage.ChatStorage
}

func NewChatHandler(storage *storage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}

func (handler ChatHandler) GetAll(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		logger.Logger.Error("failed to parse user id", slog.Any("user_id", c.Locals("user_id")))
		return errors.New("failed to parse user id")
	}
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	chats, err := handler.storage.GetAllChats()
	if err != nil {
		l.Error("failed to get list of chats", slog.String("err", err.Error()))
		return err
	}

	if len(chats) < 1 {
		c.Status(fiber.StatusNoContent)
	} else {
		c.Status(fiber.StatusOK)
	}

	l.Info("list of chats has been sent to the user", slog.Any("chats", chats))
	return c.JSON(chats)
}

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

func (handler ChatHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		logger.Logger.Error("failed to parse user id", slog.Any("user_id", c.Locals("user_id")))
		return errors.New("failed to parse user id")
	}
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	body := domain.CreateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		l.Error(
			"failed to parse request body",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	chat := models.Chat{Name: body.Name}
	if err := handler.storage.CreateChat(&chat); err != nil {
		l.Error("failed to create new chat", slog.String("err", err.Error()))
		return err
	}

	l.Info("new chat has been created", slog.Any("chat", chat))
	return c.SendStatus(fiber.StatusOK)
}

func (handler ChatHandler) Update(c *fiber.Ctx) error {
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

	body := domain.UpdateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		l.Error(
			"failed to parse request body",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.UpdateChat(uint(chatId), updates); err != nil {
		l.Error(
			"failed to update the chat",
			slog.String("err", err.Error()),
			slog.Any("updates", updates),
		)
		return err
	}

	l.Info("chat has been updated", slog.Any("updates", updates))
	return c.SendStatus(fiber.StatusOK)
}

func (handler ChatHandler) Delete(c *fiber.Ctx) error {
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

	if err := handler.storage.DeleteChat(uint(chatId)); err != nil {
		l.Error("failed to delete the chat", slog.String("err", err.Error()))
		return err
	}

	l.Info("chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}

var chatIdsToConnections = make(map[uint][]*websocket.Conn)

func (handler ChatHandler) Enter(c *websocket.Conn) {
	defer c.Close()

	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		logger.Logger.Error("failed to parse user id", slog.Any("user_id", c.Locals("user_id")))
		return
	}
	l := logger.Logger.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		l.Error("failed to parse chat id", slog.Any("chat_id", c.Params("chat_id")))
		return
	}
	l = l.With(slog.Uint64("chat_id", chatId))

	if isExists := handler.storage.IsChatExists(uint(chatId)); !isExists {
		l.Error("failed to find the chat")
		return
	}

	chatIdsToConnections[uint(chatId)] = append(chatIdsToConnections[uint(chatId)], c)
	l.Info("created new websocket connection")

	messages, err := handler.storage.GetAllMessages(uint(chatId))
	if err != nil {
		l.Error("failed to get chat history", slog.String("err", err.Error()))
		return
	}

	for _, message := range messages {
		if err := c.WriteJSON(message); err != nil {
			l.Error(
				"failed to send a message to the user",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}
	}
	l.Info("chat history has been sent to the user", slog.Any("messages", messages))

	for {
		_, text, err := c.ReadMessage()
		if err != nil {
			l.Error("failed to read the message", slog.String("err", err.Error()))
			return
		}

		message := models.Message{
			Text:   string(text),
			UserID: userId,
			ChatID: uint(chatId),
			SentAt: time.Now(),
		}
		if err := handler.storage.CreateMessage(&message); err != nil {
			l.Error(
				"failed to store the message in the database",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		} else {
			l.Info("new message has been created", slog.Any("message", message))
		}

		for _, ws := range chatIdsToConnections[uint(chatId)] {
			if ws != c {
				if err := ws.WriteJSON(message); err != nil {
					l.Error(
						"failed to send a message to other users",
						slog.String("err", err.Error()),
						slog.Any("sender_connection", c),
						slog.Any("recipient_connection", ws),
					)
					return
				}
			}
		}
		l.Info("user's message has been sent to other users", slog.Any("message", message))
	}
}
