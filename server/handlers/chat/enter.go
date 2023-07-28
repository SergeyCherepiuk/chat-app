package chathandler

import (
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/exp/slog"
)

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
	l.Info("user has been connected to the chat")

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
		l.Info("user has sent the message", slog.Any("message", message))
	}
}
