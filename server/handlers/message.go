package handlers

import (
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	storage *storage.MessageStorage
}

func NewMessageHandler(storage *storage.MessageStorage) *MessageHandler {
	return &MessageHandler{storage: storage}
}

func (handler MessageHandler) Create(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		return err
	}

	body := domain.CreateMessageRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if err := body.Validate(); err != nil {
		return err
	}

	message := models.Message{
		Text:   body.Text,
		SentAt: time.Now(),
		UserID: userId,
		ChatID: uint(chatId),
	}

	if err := handler.storage.Create(message); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
