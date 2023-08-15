package messaging

import (
	"context"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/exp/slog"
)

type MessageSenderService[Message domain.DirectMessage | domain.GroupMessage] struct{}

func NewMessageSenderService[Message domain.DirectMessage | domain.GroupMessage]() *MessageSenderService[Message] {
	return &MessageSenderService[Message]{}
}

func (sender MessageSenderService[Message]) Send(ctx context.Context, messages []Message, to ...*websocket.Conn) {
	log, loggerOk := ctx.Value(logger.LogContextKey).(logger.Logger)

	send := func(ctx context.Context, conn *websocket.Conn) {
		for _, message := range messages {
			select {
			case <-ctx.Done():
				return
			default:
				if err := conn.WriteJSON(message); err != nil && loggerOk {
					log.Error(
						"failed to send a message to the companion",
						slog.String("err", err.Error()),
						slog.Any("message", message),
					)
				}
			}
		}
	}

	for _, conn := range to {
		go send(ctx, conn)
	}
}
