package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	storage *storage.ChatStorage
}

func NewChatHandler(storage *storage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}

func (handler ChatHandler) GetAll(c *fiber.Ctx) error {
	chats, err := handler.storage.GetAllChats()
	if err != nil {
		return err
	}

	if len(chats) < 1 {
		c.Status(fiber.StatusNoContent)
	} else {
		c.Status(fiber.StatusOK)
	}
	return c.JSON(chats)
}

func (handler ChatHandler) GetById(c *fiber.Ctx) error {
	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		return err
	}

	chat, err := handler.storage.GetChatById(uint(chatId))
	if err != nil {
		return err
	}

	return c.JSON(chat)
}

func (handler ChatHandler) Create(c *fiber.Ctx) error {
	body := domain.CreateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	chat := models.Chat{Name: body.Name}
	if err := handler.storage.CreateChat(&chat); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (handler ChatHandler) Update(c *fiber.Ctx) error {
	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		return err
	}

	body := domain.UpdateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.UpdateChat(uint(chatId), updates); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (handler ChatHandler) Delete(c *fiber.Ctx) error {
	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		return err
	}

	if err := handler.storage.DeleteChat(uint(chatId)); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

var chatIdsToConnections = make(map[uint][]*websocket.Conn)

func (handler ChatHandler) Enter(c *websocket.Conn) {
	defer c.Close()

	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Println("Unathorized")
		return
	}

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		log.Printf("Invalid chat id: %s\n", err.Error())
		return
	}

	if isExists := handler.storage.IsChatExists(uint(chatId)); !isExists {
		log.Printf("Chat not found: %s\n", err.Error())
		return
	}

	chatIdsToConnections[uint(chatId)] = append(chatIdsToConnections[uint(chatId)], c)

	messages, err := handler.storage.GetAllMessages(uint(chatId))
	if err != nil {
		log.Printf("Failed to get the chat history from the database: %s\n", err.Error())
		return
	}

	for _, message := range messages {
		if err := c.WriteJSON(message); err != nil {
			log.Printf("Failed to send the chat history to the client: %s\n", err.Error())
			return
		}
	}

	for {
		_, text, err := c.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from the client: %s\n", err.Error())
			return
		}

		message := models.Message{
			Text:   string(text),
			UserID: userId,
			ChatID: uint(chatId),
			SentAt: time.Now(),
		}
		if err := handler.storage.CreateMessage(&message); err != nil {
			log.Printf("Failed to save message to the database: %s\n", err.Error())
			return
		}

		for _, ws := range chatIdsToConnections[uint(chatId)] {
			if ws != c {
				if err := ws.WriteJSON(message); err != nil {
					log.Printf("Failed to send message to other chatters: %s\n", err.Error())
					return
				}
			}
		}
	}
}
