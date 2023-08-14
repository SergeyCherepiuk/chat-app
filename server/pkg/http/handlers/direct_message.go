package handlers

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"golang.org/x/exp/slog"
)

type DirectMessageHandler struct {
	directMessageService     domain.DirectMessageService
	connectionManagerService domain.ConnectionManagerService[[2]uint]
	userService              domain.UserService
}

func NewDirectMessageHandler(
	directMessageService domain.DirectMessageService,
	connectionManagerService domain.ConnectionManagerService[[2]uint],
	userService domain.UserService,
) *DirectMessageHandler {
	return &DirectMessageHandler{
		directMessageService:     directMessageService,
		connectionManagerService: connectionManagerService,
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

	log := log.Logger{}

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

	history, err := handler.directMessageService.GetHistory(userId, companionId)
	if err != nil {
		log.Error("failed to get chat history", slog.String("err", err.Error()))
		return
	}

	for _, message := range history {
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

		for _, ws := range handler.connectionManagerService.GetConnections(key).Values() {
			if ws != c {
				if err := ws.(*websocket.Conn).WriteJSON(message); err != nil {
					log.Error(
						"failed to send a message to the companion",
						slog.String("err", err.Error()),
						slog.Any("message", message),
					)
					return
				}
			}
		}
		log.Info("user has sent the message", slog.Any("message", message))
	}
}

func (handler DirectMessageHandler) UpdateMessage(c *fiber.Ctx) error {
	log := log.Logger{}

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
	log := log.Logger{}

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
	log := log.Logger{}

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
