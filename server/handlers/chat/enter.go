package chathandler

import (
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/exp/slog"
)

var chatIdsToConnections = make(map[uint]*hashset.Set)

func (handler ChatHandler) Enter(c *websocket.Conn) {
	defer c.Close()

	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		log.Error(
			"failed to parse chat id",
			slog.String("err", err.Error()),
			slog.Any("chat_id", c.Params("chat_id")),
		)
		return
	}
	log.With(slog.Uint64("chat_id", chatId))

	if isExists := handler.storage.IsChatExists(uint(chatId)); !isExists {
		log.Error("failed to find the chat", slog.String("err", err.Error()))
		return
	}

	if _, ok := chatIdsToConnections[uint(chatId)]; !ok {
		chatIdsToConnections[uint(chatId)] = hashset.New()
	}
	chatIdsToConnections[uint(chatId)].Add(c)
	log.Info("user has been connected to the chat")
	defer func() {
		chatIdsToConnections[uint(chatId)].Remove(c)
	}()

	messages, err := handler.storage.GetAllMessages(uint(chatId))
	if err != nil {
		log.Error("failed to get chat history", slog.String("err", err.Error()))
		return
	}

	for _, message := range messages {
		if err := c.WriteJSON(message); err != nil {
			log.Error(
				"failed to sent the message to the user",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}
	}

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

		message := models.Message{
			Text:   string(text),
			UserID: userId,
			ChatID: uint(chatId),
			SentAt: time.Now(),
		}
		if err := handler.storage.CreateMessage(&message); err != nil {
			log.Error(
				"failed to store the message",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}

		for _, ws := range chatIdsToConnections[uint(chatId)].Values() {
			if ws != c {
				if err := ws.(*websocket.Conn).WriteJSON(message); err != nil {
					log.Error(
						"failed to send a message to other users",
						slog.String("err", err.Error()),
						slog.Any("sender_connection", c),
						slog.Any("recipient_connection", ws),
					)
					return
				}
			}
		}
		log.Info("user has sent the message", slog.Any("message", message))
	}
}
