package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/contrib/websocket"
)

type ChatHandler struct {
	storage *storage.ChatStorage
}

func NewChatHandler(storage *storage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}

var connections []*websocket.Conn

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

	if isExists := handler.storage.IsExists(uint(chatId)); !isExists {
		log.Printf("Chat not found: %s\n", err.Error())
		return
	}

	connections = append(connections, c)

	messages, err := handler.storage.GetAllForChat(uint(chatId))
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
		if err := handler.storage.Create(&message); err != nil {
			log.Printf("Failed to save message to the database: %s\n", err.Error())
			return
		}

		for _, ws := range connections {
			if ws != c {
				if err := ws.WriteJSON(message); err != nil {
					log.Printf("Failed to send message to other chatters: %s\n", err.Error())
					return
				}
			}
		}
	}
}
