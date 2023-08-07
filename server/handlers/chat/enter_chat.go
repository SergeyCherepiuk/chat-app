package chathandler

import (
	chatdomain "github.com/SergeyCherepiuk/chat-app/domain/chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/exp/slog"
)

type pair struct {
	first  uint
	second uint
}

func (pair pair) getKey() [2]uint {
	if pair.first > pair.second {
		return [2]uint{pair.second, pair.first}
	}
	return [2]uint{pair.first, pair.second}
}

var connections = make(map[[2]uint]*hashset.Set)

func (handler ChatHandler) EnterChat(c *websocket.Conn) {
	defer c.Close()

	log := logger.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	companionId := c.Locals("companion_id").(uint)
	log.With(slog.Uint64("companion_id", uint64(companionId)))

	key := pair{first: userId, second: companionId}.getKey()
	if _, ok := connections[key]; !ok {
		connections[key] = hashset.New()
	}
	connections[key].Add(c)
	log.Info("user has been connected to the chat", slog.Any("key", key))
	defer func() {
		connections[key].Remove(c)
	}()

	messages, err := handler.chatStorage.GetHistory(userId, companionId)
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

		body := chatdomain.CreateMessageBody{Message: string(text)}
		if err := body.Validate(); err != nil {
			log.Error(
				"body isn't valid",
				slog.String("err", err.Error()),
				slog.Any("body", body),
			)
			continue
		}

		message := models.ChatMessage{
			Message:  body.Message,
			From:     userId,
			To:       companionId,
			IsEdited: false,
		}
		if err := handler.chatStorage.Create(&message); err != nil {
			log.Error(
				"failed to store the message",
				slog.String("err", err.Error()),
				slog.Any("message", message),
			)
			return
		}

		for _, ws := range connections[key].Values() {
			if ws != c {
				if err := ws.(*websocket.Conn).WriteJSON(message); err != nil {
					log.Error("failed to send a message to the companion", slog.String("err", err.Error()))
					return
				}
			}
		}
		log.Info("user has sent the message", slog.Any("message", message))
	}
}
