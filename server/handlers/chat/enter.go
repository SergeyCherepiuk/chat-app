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

	userId, _ := c.Locals("user_id").(uint)

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse chat id",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Any("chat_id", c.Params("chat_id")),
			},
		}
		return
	}

	if isExists := handler.storage.IsChatExists(uint(chatId)); !isExists {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to find the chat",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Uint64("chat_id", chatId),
			},
		}
		return
	}

	if _, ok := chatIdsToConnections[uint(chatId)]; !ok {
		chatIdsToConnections[uint(chatId)] = hashset.New()
	}
	chatIdsToConnections[uint(chatId)].Add(c)
	logger.LogMessages <- logger.LogMessage{
		Message: "user has been connected to the chat",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(userId)),
			slog.Uint64("chat_id", chatId),
		},
	}
	defer func() {
		chatIdsToConnections[uint(chatId)].Remove(c)
	}()

	messages, err := handler.storage.GetAllMessages(uint(chatId))
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to get chat history",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
				slog.Uint64("chat_id", chatId),
			},
		}
		return
	}

	for _, message := range messages {
		if err := c.WriteJSON(message); err != nil {
			logger.LogMessages <- logger.LogMessage{
				Message: "failed to get chat history",
				Level:   slog.LevelError,
				Attrs: []slog.Attr{
					slog.String("err", err.Error()),
					slog.Uint64("user_id", uint64(userId)),
					slog.Uint64("chat_id", chatId),
					slog.Any("message", message),
				},
			}
			return
		}
	}

	for {
		_, text, err := c.ReadMessage()
		if websocket.IsCloseError(err, 1000, 1005) {
			logger.LogMessages <- logger.LogMessage{
				Message: "user has been disconnected",
				Level:   slog.LevelInfo,
				Attrs: []slog.Attr{
					slog.String("err", err.Error()),
					slog.Uint64("user_id", uint64(userId)),
					slog.Uint64("chat_id", chatId),
				},
			}
			return
		}
		if err != nil {
			logger.LogMessages <- logger.LogMessage{
				Message: "failed to read the message",
				Level:   slog.LevelError,
				Attrs: []slog.Attr{
					slog.String("err", err.Error()),
					slog.Uint64("user_id", uint64(userId)),
					slog.Uint64("chat_id", chatId),
				},
			}
			return
		}

		message := models.Message{
			Text:   string(text),
			UserID: userId,
			ChatID: uint(chatId),
			SentAt: time.Now(),
		}
		if err := handler.storage.CreateMessage(&message); err != nil {
			logger.LogMessages <- logger.LogMessage{
				Message: "failed to store the message in the database",
				Level:   slog.LevelError,
				Attrs: []slog.Attr{
					slog.String("err", err.Error()),
					slog.Uint64("user_id", uint64(userId)),
					slog.Uint64("chat_id", chatId),
					slog.Any("message", message),
				},
			}
			return
		}

		for _, ws := range chatIdsToConnections[uint(chatId)].Values() {
			if ws != c {
				if err := ws.(*websocket.Conn).WriteJSON(message); err != nil {
					logger.LogMessages <- logger.LogMessage{
						Message: "failed to send a message to other users",
						Level:   slog.LevelError,
						Attrs: []slog.Attr{
							slog.String("err", err.Error()),
							slog.Uint64("user_id", uint64(userId)),
							slog.Uint64("chat_id", chatId),
							slog.Any("sender_connection", c),
							slog.Any("recipient_connection", ws),
						},
					}
					return
				}
			}
		}
		logger.LogMessages <- logger.LogMessage{
			Message: "user has sent the message",
			Level:   slog.LevelInfo,
			Attrs: []slog.Attr{
				slog.Uint64("user_id", uint64(userId)),
				slog.Uint64("chat_id", chatId),
				slog.Any("message", message),
			},
		}
	}
}
