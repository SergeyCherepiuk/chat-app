package messaging

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/exp/slog"
)

type MessageSenderService[Message domain.DirectMessage | domain.GroupMessage] struct{}

func NewMessageSenderService[Message domain.DirectMessage | domain.GroupMessage]() *MessageSenderService[Message] {
	return &MessageSenderService[Message]{}
}

func (sender MessageSenderService[Message]) Send(messages []Message, to ...*websocket.Conn) {
	log := log.Logger{}

	send := func(conn *websocket.Conn) {
		for _, message := range messages {
			if err := conn.WriteJSON(message); err != nil {
				// TODO: Figure out how to get other log attrs (user_id and chat_id)
				log.Error(
					"failed to send a message to the companion",
					slog.String("err", err.Error()),
					slog.Any("message", message),
				)
			}
		}
	}

	for _, conn := range to {
		go send(conn)
	}
}
